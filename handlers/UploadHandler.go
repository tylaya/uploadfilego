package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// UploadFile ini fungsi for mo handle proses upload
func UploadFile(w http.ResponseWriter, r *http.Request) {
	// 1. Kalo pake method GET, dia mo muncul tu form HTML #eaaaa
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<!DOCTYPE html>
			<html lang="id">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Tugas Upload File</title>
				<style>
					body {
						font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
						background-color: #f4f7f6;
						display: flex;
						justify-content: center;
						align-items: center;
						height: 100vh;
						margin: 0;
					}
					.upload-card {
						background-color: #ffffff;
						padding: 40px;
						border-radius: 12px;
						box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
						text-align: center;
						width: 100%;
						max-width: 400px;
					}
					.upload-card h2 {
						color: #333333;
						margin-top: 0;
						margin-bottom: 10px;
					}
					.upload-card p {
						color: #666666;
						font-size: 14px;
						margin-bottom: 25px;
					}
					.file-input {
						width: 90%;
						padding: 10px;
						border: 2px dashed #cccccc;
						border-radius: 8px;
						background-color: #f9f9f9;
						cursor: pointer;
						margin-bottom: 20px;
						transition: border-color 0.3s;
					}
					.file-input:hover {
						border-color: #00ADD8; /* Warna khas Golang */
					}
					.btn-submit {
						background-color: #00ADD8;
						color: white;
						border: none;
						padding: 12px 20px;
						font-size: 16px;
						font-weight: bold;
						border-radius: 8px;
						cursor: pointer;
						width: 100%;
						transition: background-color 0.3s;
					}
					.btn-submit:hover {
						background-color: #008db3;
					}
				</style>
			</head>
			<body>
				<div class="upload-card">
					<h2>Tugas Mandiri Golang</h2>
					<p>Silakan pilih file dari perangkat Anda untuk diunggah ke server.</p>
					<form action="/upload" method="post" enctype="multipart/form-data">
						<input type="file" name="file_tugas" class="file-input" required />
						<button type="submit" class="btn-submit">Unggah File</button>
					</form>
				</div>
			</body>
			</html>
		`)
		return
	}

	// 2. Kalo method POST, itu dia mba proses dp file yang torang unggah
	if r.Method == http.MethodPost {
		// Ini mo kse batas dp ukuran file maksimal 10 MB biar apa? biarin #izinnn
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Ukuran file terlalu besar", http.StatusBadRequest)
			return
		}

		// pokoknya ini for ambe tu file "file_tugas"
		file, handler, err := r.FormFile("file_tugas")
		if err != nil {
			http.Error(w, "Gagal mengambil file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// --- Checklist 1: Batasi tipe file ---
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			http.Error(w, "Cannot read file", http.StatusInternalServerError)
			return
		}

		// Kembalikan pointer file ke awal setelah membaca 512 byte, biar filenya gak corrupt pas disalin!
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			http.Error(w, "Error saat memproses file", http.StatusInternalServerError)
			return
		}

		filetype := http.DetectContentType(buffer)
		fmt.Println("Type:", filetype)

		allowed := map[string]bool{
			"image/jpeg":      true,
			"image/png":       true,
			"application/pdf": true,
		}

		if !allowed[filetype] {
			// Set status code menjadi 400 (Bad Request)
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "text/html")

			// Tampilkan halaman HTML Error
			fmt.Fprint(w, `
				<!DOCTYPE html>
				<html lang="id">
				<head>
					<meta charset="UTF-8">
					<meta name="viewport" content="width=device-width, initial-scale=1.0">
					<title>Upload Gagal</title>
					<style>
						body { font-family: 'Segoe UI', sans-serif; background-color: #f4f7f6; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
						.error-card { background-color: #ffffff; padding: 40px; border-radius: 12px; box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1); text-align: center; max-width: 400px; width: 100%; border-top: 5px solid #dc3545; }
						h2 { color: #dc3545; margin-top: 0; }
						p { color: #555; line-height: 1.5; margin-bottom: 25px; }
						.btn-back { display: inline-block; background-color: #6c757d; color: white; text-decoration: none; padding: 10px 20px; border-radius: 8px; font-weight: bold; transition: background-color 0.3s; }
						.btn-back:hover { background-color: #5a6268; }
					</style>
				</head>
				<body>
					<div class="error-card">
						<h2>⚠️ Upload Gagal!</h2>
						<p><b>Invalid file type.</b><br>Sistem hanya menerima file dengan format <b>JPEG, PNG, atau PDF</b>.</p>
						<a href="/upload" class="btn-back">Coba Lagi</a>
					</div>
				</body>
				</html>
			`)
			return
		}

		// Checklist 2: Batasi penamaan file
		// Ambil ekstensi dari file asli (misal: .jpg, .pdf)
		ext := filepath.Ext(handler.Filename)

		// Tambahkan jam dan menit agar kalau upload di hari yang sama, namanya tidak bentrok
		now := time.Now().Format("2006-01-02_15-04-05")
		filename := now + ext

		fmt.Println("Uploaded Original Filename:", handler.Filename)
		fmt.Println("Saved As:", filename)
		fmt.Println("File Size:", handler.Size)
		fmt.Println("MIME Header:", handler.Header)

		// Checklist 3: Simpan File
		// Bekeng folder "uploads" klo blum adaa, klo so ada yasudah
		uploadDir := "uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			http.Error(w, "Gagal membuat folder penyimpanan", http.StatusInternalServerError)
			return
		}

		// Buat file tujuan di folder uploads
		filePath := filepath.Join(uploadDir, filename)
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy isi file yang diunggah ke file tujuan
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		// Kase tau kalo gacor or sukses
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html lang="id">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Upload Berhasil</title>
				<style>
					body { font-family: 'Segoe UI', sans-serif; background-color: #f4f7f6; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
					.success-card { background-color: #ffffff; padding: 40px; border-radius: 12px; box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1); text-align: center; max-width: 400px; width: 100%%; }
					h2 { color: #28a745; margin-top: 0; }
					p { color: #555; line-height: 1.5; margin-bottom: 25px; }
					.btn-back { display: inline-block; background-color: #6c757d; color: white; text-decoration: none; padding: 10px 20px; border-radius: 8px; font-weight: bold; transition: background-color 0.3s; }
					.btn-back:hover { background-color: #5a6268; }
				</style>
			</head>
			<body>
				<div class="success-card">
					<h2>🎉 Upload Berhasil!</h2>
					<p>File asli <b>%s</b> telah berhasil disimpan dengan nama baru <b>%s</b> ke dalam folder server.</p>
					<a href="/upload" class="btn-back">Kembali Upload</a>
				</div>
			</body>
			</html>
		`, handler.Filename, filename)
		return
	}

	// Slain ini GET or POST yaa goodbyee sayonara karnah nda di #izinnn
	http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
}
