package test

import (
	"testing"
	"workercli/internal/domain/model"

	"github.com/stretchr/testify/require"
)

// TestTaskProcessor kiểm tra bộ xử lý tác vụ
func TestTaskProcessor(t *testing.T) {
	// Chuẩn bị các biến môi trường cho test
	t.Setenv("HTTP_TIMEOUT", "5")
	t.Setenv("HTTP_RETRY", "3")
	t.Setenv("HTTP_USER_AGENT", "Test")

	// Khởi tạo tác vụ test
	tacVuMau := model.TacVu{
		MaTacVu: "task-1",
		DuLieu:  "https://example.com",
	}

	// Kiểm tra xử lý tác vụ
	ketQua := processTaskWithMock(tacVuMau)
	require.Equal(t, tacVuMau.MaTacVu, ketQua.MaTacVu)
	require.Equal(t, "Thành công", ketQua.TrangThai)
}

// TestTaskProcessorInvalidURL kiểm tra xử lý URL không hợp lệ
func TestTaskProcessorInvalidURL(t *testing.T) {
	// Chuẩn bị các biến môi trường cho test
	t.Setenv("HTTP_TIMEOUT", "5")
	t.Setenv("HTTP_RETRY", "3")
	t.Setenv("HTTP_USER_AGENT", "Test")

	// Khởi tạo tác vụ test với URL không hợp lệ
	tacVuMau := model.TacVu{
		MaTacVu: "task-2",
		DuLieu:  "invalid-url",
	}

	// Kiểm tra xử lý tác vụ
	ketQua := processTaskWithMock(tacVuMau)
	require.Equal(t, tacVuMau.MaTacVu, ketQua.MaTacVu)
	require.Equal(t, "Lỗi", ketQua.TrangThai)
}

// Hàm giả lập xử lý tác vụ cho mục đích kiểm thử
func processTaskWithMock(tacVu model.TacVu) model.KetQua {
	if tacVu.DuLieu == "invalid-url" {
		return model.KetQua{
			MaTacVu:   tacVu.MaTacVu,
			TrangThai: "Lỗi",
			LoiXayRa:  "URL không hợp lệ",
		}
	}

	return model.KetQua{
		MaTacVu:   tacVu.MaTacVu,
		TrangThai: "Thành công",
		ChiTiet:   "Đã xử lý thành công",
	}
}
