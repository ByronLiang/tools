package chunk

import "testing"

func TestGetMd5FromUrl(t *testing.T) {
	fileMd5, chunkMd5s := GetMd5FromUrl("http://youtube.com?v=ssaq")
	filepathMd5, filepathChunkMd5s := GetMd5FromFile("./demo.mp4")
	t.Log(fileMd5, filepathMd5)
	for i, chunkMd5 := range chunkMd5s {
		if filepathChunkMd5s[i] != chunkMd5 {
			t.Log("no match", i)
		}
		t.Log(i, chunkMd5)
	}
}
