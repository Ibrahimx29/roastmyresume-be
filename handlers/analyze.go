package handlers

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"roastmyresume/utils"

	"github.com/gin-gonic/gin"
)

type AnalyzeRequest struct {
    Mode string `form:"mode"` // "roast" or "serious"
}

func AnalyzeResume(c *gin.Context) {
    // ✅ Parse form data including 'mode'
    var req AnalyzeRequest
    if err := c.ShouldBind(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing mode"})
        return
    }

    // ✅ Handle uploaded file
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file"})
        return
    }

    // ✅ Save file temporarily
    tempDir := os.TempDir()
    tempFilePath := filepath.Join(tempDir, file.Filename)
    out, err := os.Create(tempFilePath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create temp file"})
        return
    }
    defer out.Close()

    uploaded, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open uploaded file"})
        return
    }
    defer uploaded.Close()

    _, err = io.Copy(out, uploaded)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
        return
    }

    // ✅ Call Python script
    cmd := exec.Command("python3", "./extract_text.py", tempFilePath)
    output, err := cmd.CombinedOutput()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Python error: " + err.Error(), "details": string(output)})
        return
    }

    text := string(output)
    prompt := utils.GeneratePrompt(text, req.Mode)
    response, err := utils.CallGroq(prompt)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"feedback": response, "resume_text": text})
}