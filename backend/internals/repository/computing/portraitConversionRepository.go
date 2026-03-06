package computingRepository

import (
	"fmt"
	"os/exec"
	"bytes"
	"strconv"
)

type PortraitConversionRepository struct {
	webpImageQuality int
}

func NewPortraitConversionRepository(webpImageQuality int) *PortraitConversionRepository {
	return &PortraitConversionRepository{
		webpImageQuality: webpImageQuality,
	}
}

func (r *PortraitConversionRepository) ConvertToWebP(input []byte) ([]byte, error) {
    cmd := exec.Command("ffmpeg",
        "-i", "pipe:0",
        "-quality", strconv.Itoa(r.webpImageQuality),
        "-f", "webp",
        "pipe:1",
    )

    cmd.Stdin = bytes.NewReader(input)

    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    if err := cmd.Run(); err != nil {
        return nil, fmt.Errorf("ffmpeg: %w: %s", err, stderr.String())
    }

    return stdout.Bytes(), nil
}