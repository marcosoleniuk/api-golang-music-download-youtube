package services

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"youtube_converter/common"
)

const CookieFile = "./cookies.txt"

func DownloadVideo(link, format, output, taskID string) error {
	if err := validateCookies(); err != nil {
		return err
	}

	args := buildArgs(format, output, link)
	cmd := exec.Command("yt-dlp", args...)
	log.Printf("Executing command: %s", cmd.String())

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	go handleOutput(stdout, taskID)

	return cmd.Wait()
}

func validateCookies() error {
	if _, err := os.Stat(CookieFile); os.IsNotExist(err) {
		log.Printf("Error: cookie file not found!")
		return fmt.Errorf("cookie file not found: %s", CookieFile)
	}
	return nil
}

func buildArgs(format, output, link string) []string {
	args := []string{
		"-o", output,
		"--no-playlist",
		"--cookies", CookieFile,
		"--concurrent-fragments", "32",
		"--progress",
		"--newline",
	}

	if format == "mp3" {
		args = append(args, "--extract-audio", "--audio-format", "mp3", "--audio-quality", "0")
	} else {
		args = append(args, "--format", fmt.Sprintf("bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best"))
	}

	args = append(args, link)
	return args
}

func handleOutput(stdout io.ReadCloser, taskID string) {
	defer func(stdout io.ReadCloser) {
		err := stdout.Close()
		if err != nil {
			log.Printf("Error closing stdout: %v", err)
		}
	}(stdout)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "[download]") {
			progress := parseProgress(line)
			common.Broadcast <- common.ProgressMessage{TaskID: taskID, Percentage: progress}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from stdout: %v", err)
	}
}

func parseProgress(line string) float64 {
	re := regexp.MustCompile(`(\d+\.\d+)%`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		progress, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			return progress
		}
		log.Printf("Error parsing progress: %v", err)
	}
	return 0.0
}
