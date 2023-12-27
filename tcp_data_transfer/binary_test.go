package tcp_data_transfer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"testing"
)

func TestBinary(t *testing.T) {
	// 정수 데이터
	num := uint16(42)

	// 파일 열기
	file, err := os.Create("binary_data.bin") // create binary file
	if err != nil {
		fmt.Println("파일을 열 수 없습니다.", err)
		return
	}
	defer file.Close()

	// 정수를 바이너리로 인코딩하여 파일에 쓰기
	err = binary.Write(file, binary.LittleEndian, num) // little endian 작은 자릿수부터 바이트 저장 <-> big endian 은 반대로 큰 자릿수부터 바이트 저장
	if err != nil {
		fmt.Println("파일에 데이터를 쓸 수 없습니다.", err)
		return
	}

	fmt.Println("파일에 데이터를 성공적으로 썼습니다.")
}

func TestBinaryMapping(t *testing.T) {
	data := uint16(500)
	buf := new(bytes.Buffer)

	// 리틀 엔디안으로 데이터를 바이너리로 인코딩
	binary.Write(buf, binary.LittleEndian, data)

	// 결과 출력
	fmt.Printf("Encoded data: %v\n", buf.Bytes())

	type DataStruct struct {
		Field1 uint16
		Field2 uint32
		Field3 uint16
	}

	st := DataStruct{Field1: 300, Field2: 200, Field3: 256}
	buf = new(bytes.Buffer)

	// 구조체를 바이너리로 인코딩
	binary.Write(buf, binary.LittleEndian, st)

	// 결과 출력
	fmt.Printf("Encoded data: %v\n", buf.Bytes())
}

func TestBuffer(t *testing.T) {
	var buf bytes.Buffer

	buf.Write([]byte("안녕! 나는"))
	buf.Write([]byte("World!!"))

	data := make([]byte, 20)
	n, _ := buf.Read(data)

	fmt.Printf("Read %d bytes: %s\n", n, data[:n])
}
