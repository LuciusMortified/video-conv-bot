package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/uuid"

	"github.com/LuciusMortified/video-conv-bot/internal/ent"
	"github.com/LuciusMortified/video-conv-bot/pkg/logger"
	"github.com/LuciusMortified/video-conv-bot/pkg/ptr"
)

type Service struct {
	client http.Client
	cfg    Config
}

type Config struct {
	StoragePath string
}

func (c Config) Validate() error {
	if c.StoragePath == "" {
		return fmt.Errorf("storage path is required")
	}
	return nil
}

func New(cfg Config) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("cfg.Validate: %w", err)
	}
	return &Service{
		client: http.Client{},
		cfg:    cfg,
	}, nil
}

func (s *Service) Convert(ctx context.Context, params ent.ConvertParams) ent.ConvertStateChan {
	stateChan := make(ent.ConvertStateChan)

	// TODO: limit max converting in one time
	go s.doConvert(ctx, params, stateChan)

	return stateChan
}

func (s *Service) doConvert(ctx context.Context, params ent.ConvertParams, stateChan ent.ConvertStateChan) {
	sp := &statePusher{ch: stateChan}
	defer sp.Close()

	origFilename, err := s.downloadFile(ctx, sp, params.URL)
	if err != nil {
		logger.With("error", err).Error("Failed to download file")
		return
	}
	defer func() {
		_ = os.Remove(origFilename)
	}()

	newFilename, err := s.convertFile(ctx, sp, origFilename)
	if err != nil {
		logger.With("error", err).Error("Failed to convert file")
		return
	}

	data, err := s.readFile(ctx, sp, newFilename)
	if err != nil {
		logger.With("error", err).Error("Failed to read file")
		return
	}

	sp.Done(newFilename, data)
}

func (s *Service) downloadFile(ctx context.Context, sp *statePusher, rawURL string) (string, error) {
	sp.Downloading()

	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		sp.Error("Не удалось скачать файл")
		return "", fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	logger.With("url", rawURL).Info("Download file")
	resp, err := s.client.Do(req)
	if err != nil {
		sp.Error("Не удалось скачать файл")
		return "", fmt.Errorf("s.client.Do: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	mediatype, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		sp.Error("Не удалось проверить формат файла")
		return "", fmt.Errorf("mime.ParseMediaType: %w", err)
	}

	ok := strings.Contains(mediatype, "video/") ||
		strings.Contains(mediatype, "application/octet-stream")
	if !ok {
		sp.Error("Неподдерживаемый формат")
		return "", errors.New("unsupported file format")
	}

	name := uuid.New().String()
	origFilename := path.Join(s.cfg.StoragePath, name)

	f, err := os.OpenFile(origFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		sp.Error("Не удалось скачать файл")
		return "", fmt.Errorf("os.OpenFile: %w", err)
	}
	defer func(f *os.File) {
		_ = f.Close()
		if err != nil {
			_ = os.Remove(origFilename)
		}
	}(f)

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		sp.Error("Не удалось скачать файл")
		return "", fmt.Errorf("io.Copy: %w", err)
	}

	return origFilename, nil
}

func (s *Service) convertFile(_ context.Context, sp *statePusher, origFilename string) (string, error) {
	sp.Converting()

	name := fmt.Sprintf("%s.mp4", uuid.New().String())
	newFilename := path.Join(s.cfg.StoragePath, name)

	cmd := exec.Command("ffmpeg", "-i", origFilename, newFilename)
	if err := cmd.Run(); err != nil {
		sp.Error("Не удалось конвертировать файл")
		return "", fmt.Errorf("cmd.Run: %w", err)
	}

	return newFilename, nil
}

func (s *Service) readFile(_ context.Context, sp *statePusher, newFilename string) (io.ReadCloser, error) {
	f, err := os.OpenFile(newFilename, os.O_RDONLY, 0600)
	if err != nil {
		sp.Error("Не удалось прочитать новый файл")
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}

	return &deleteOnCloseFile{filename: newFilename, file: f}, nil
}

type deleteOnCloseFile struct {
	filename string
	file     *os.File
}

func (f *deleteOnCloseFile) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func (f *deleteOnCloseFile) Close() error {
	if err := f.file.Close(); err != nil {
		return err
	}

	if err := os.Remove(f.filename); err != nil {
		return err
	}

	return nil
}

type statePusher struct {
	ch ent.ConvertStateChan
}

func (s *statePusher) Unsupported() {
	s.ch <- ent.ConvertState{
		Status: ent.ConvertUnsupported,
	}
}

func (s *statePusher) Downloading() {
	s.ch <- ent.ConvertState{
		Status: ent.ConvertDownloading,
	}
}

func (s *statePusher) Error(msg string) {
	s.ch <- ent.ConvertState{
		Status: ent.ConvertError,
		Error:  ptr.AnyRef(msg),
	}
}

func (s *statePusher) Converting() {
	s.ch <- ent.ConvertState{
		Status: ent.ConvertConverting,
	}
}

func (s *statePusher) Done(filename string, data io.ReadCloser) {
	s.ch <- ent.ConvertState{
		Status: ent.ConvertDone,
		Result: &ent.ConvertResult{
			Filename: filename,
			Data:     data,
		},
	}
}

func (s *statePusher) Close() {
	close(s.ch)
}
