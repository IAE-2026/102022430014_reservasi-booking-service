# Prompt 1
Tambahkan variabel environment baru di .env dan perbarui file konfigurasi di internal/infrastructure/. Tambahkan URL untuk SSO_URL=https://iae-sso.virtualfri.id dan API_KEY=KEY-MHS-25. Kemudian, buat sebuah helper atau service M2M di layer infrastructure untuk melakukan HTTP POST ke /api/v1/auth/token menggunakan payload JSON {"api_key": "KEY-MHS-25"}. Fungsi ini harus mengembalikan Bearer JWT, dan cache token tersebut di Redis dengan TTL menyesuaikan expired token agar sistem tidak terus-menerus memanggil API dari server cloud.


# Goal Description 1
# Integrasi Autentikasi M2M SSO Selesai

Sesuai dengan target pada `sequence.md`, saya telah membuat logika dasar yang diperlukan untuk mengintegrasikan Layanan Reservasi dengan server SSO untuk komunikasi *Machine-to-Machine* (M2M). 

Berikut adalah rangkuman dari pembaruan yang telah saya buat:

### 1. Penambahan Konfigurasi Environment
Saya telah memperbarui file `.env` dan `.env.example` untuk mencakup konfigurasi kredensial SSO yang diinstruksikan:

```env
# SSO Configuration
SSO_URL=https://iae-sso.virtualfri.id
API_KEY=KEY-MHS-25
```

### 2. Implementasi SSO M2M Helper
Saya telah menambahkan *file* baru di layer infrastructure, yaitu **`internal/infrastructure/sso.go`**. 
Fungsi utama di dalam *file* ini adalah `GetM2MToken(ctx context.Context)`, yang akan:
- **Cek Cache (Redis)**: Mencari kunci `sso:m2m_token`. Jika *cache hit*, sistem tidak akan memanggil server SSO sama sekali.
- **Fetch SSO API**: Jika *cache miss*, sistem akan menembak HTTP POST ke `${SSO_URL}/api/v1/auth/token` menggunakan body JSON `{"api_key": "..."}`.
- **Ekstraksi Waktu Expired Token (TTL)**: Secara mandiri melakukan proses *base64 decode* pada bagian *payload* JWT (bagian tengah token) untuk mengambil klaim `exp` (waktu kadaluarsa). Sistem **tidak** bergantung pada library JWT pihak ketiga untuk menghindari masalah kompatibilitas, melainkan menggunakan `base64.RawURLEncoding` standar Golang yang sangat ringan dan presisi.
- **Update Cache Redis**: Menyimpan JWT ke Redis dengan durasi *Time-To-Live* (TTL) secara persis sama dengan sisa waktu berlakunya JWT tersebut (dikurangi batas toleransi 1 menit agar lebih aman).

### Verifikasi
Kompilasi kode program menggunakan *Go compiler* (`go build`) telah **berhasil dengan sempurna** (Exit code: 0), yang menandakan tidak ada *syntax error* maupun *import cycle* pada integrasi ini. 

> [!TIP]
> **Selanjutnya**
> Fungsi `infrastructure.GetM2MToken(ctx)` kini siap digunakan oleh *Usecase Layer* Anda untuk membubuhkan *Bearer Token* saat hendak menembak integrasi SOAP API (`/soap/v1/audit`) atau Message Broker API (`/api/v1/messages/publish`).


# Prompt 2
Perbarui middleware autentikasi Gin yang sudah ada. Selain mengecek header X-IAE-KEY, tambahkan pengecekan Authorization: Bearer <token>. Buat fungsi untuk mengambil Public Key (RS256) dari GET https://iae-sso.virtualfri.id/api/v1/auth/jwks dan cache JWKS tersebut di Redis selama 24 jam. Gunakan library golang-jwt/jwt dan JWKS tersebut untuk memverifikasi keaslian JWT klien. Jika valid, ekstrak claims email-nya, lalu lakukan query menggunakan GORM ke tabel users untuk memvalidasi role lokalnya. Lempar identitas user ke Gin context.

# Goal Description 2
# Integrasi SSO (M2M & Autentikasi Klien)

Proses pengamanan Layanan Reservasi berdasarkan `sequence.md` (Fase 1: Autentikasi & Otorisasi) telah berhasil diimplementasikan secara komprehensif.

Berikut adalah fitur yang telah dirampungkan:

---

## 1. M2M Authentication Helper (`internal/infrastructure/sso.go`)
Helper ini berfungsi untuk memungkinkan layanan kita "berbicara" dengan *Cloud Dosen* melalui *Bearer Token* saat hendak mempublikasikan *Audit* atau event *RabbitMQ*.
- **Otomatisasi Caching**: Menggunakan Redis untuk men-*cache* JWT berdasarkan nilai aktual klaim `exp` yang ter-enkode (tanpa membebani aplikasi dengan *parsing* yang kompleks).
- **Graceful Fallback**: Jika Redis mati atau token tidak ditemukan, akan otomatis memanggil `/api/v1/auth/token`.

---

## 2. JWT Gin Middleware (`pkg/middleware/auth.go`)
Middleware otorisasi utama telah berevolusi dari sekadar pengecekan statis `X-IAE-KEY` menjadi sistem verifikasi berlapis berbasis token asimetris:

### A. Pengambilan & Penyimpanan Public Key (JWKS)
- Diinisiasi secara dinamis (lazy-loading) dari `https://iae-sso.virtualfri.id/api/v1/auth/jwks`.
- Disimpan *(cached)* ke dalam Redis dengan masa tenggang (TTL) 24 jam untuk memaksimalkan efisiensi (mengurangi *overhead* komunikasi antar-jaringan).

### B. Validasi JWT dengan `keyfunc`
- Menggunakan pustaka *industry standard* `github.com/golang-jwt/jwt/v5` yang dipadukan dengan `github.com/MicahParks/keyfunc/v2`.
- `keyfunc` memastikan bahwa setiap struktur JSON JWKS yang diunduh langsung dapat dipetakan secara matematis ke kunci kriptografi `rsa.PublicKey` yang sah, menghalau token tiruan.

### C. Pemetaan ke Database Lokal (Role-Based Authorization)
- Mengambil klaim `email` dari JWT setelah tanda tangannya dinyatakan valid.
- Secara transparan melakukan `SELECT id, email, role FROM users WHERE email = ?` menggunakan koneksi GORM ke PostgreSQL.
- Jika pengguna terdaftar, identitas (`ID`, `Email`, `Role`) tersebut langsung diinjeksi ke `gin.Context`.
- Identitas ini siap dibaca oleh *Delivery Handler* menggunakan `c.Get("userRole")` untuk membedakan perlakuan antara profil `guest` dan `admin`.

> [!TIP]
> **Akses Identitas di Controller**
> Untuk mendapatkan email atau role di level handler, Anda bisa menggunakan sintaks bawaan Gin:
> ```go
> role := c.GetString("userRole")
> email := c.GetString("userEmail")
> ```

### Hasil Kompilasi
Perubahan ini telah di-*build* menggunakan kompiler golang terbaru dan berjalan **100% tanpa error** sintaks maupun siklus dependensi (`exit code: 0`).


# Prompt 3
Buat file soap_client.go di internal/infrastructure/. Buat fungsi HTTP POST ke /soap/v1/audit. Gunakan M2M Token dari helper yang dibuat sebelumnya di header Authorization. Fungsi ini menerima payload JSON transaksi, lalu membungkusnya ke dalam format XML Envelope yang memiliki tag <iae:TeamID>TEAM-25</iae:TeamID>, <iae:ActivityName>BookingCreated</iae:ActivityName>, dan <iae:LogContent><![CDATA[ ...JSON... ]]></iae:LogContent>. Lakukan unmarshal pada respons XML dari server cloud untuk mengekstrak dan mengembalikan nilai dari tag <iae:ReceiptNumber>.

# Goal Description 3
# Integrasi SOAP Client

Fitur pelaporan riwayat pemesanan *(Audit Trail)* ke layanan pusat dosen via arsitektur SOAP telah sukses diimplementasikan dan dikirim ke GitHub. 

Berikut rincian dari fitur yang ditambahkan:

---

## 1. SOAP XML Enveloping (`soap_client.go`)
Fungsi `SendAuditLog` secara spesifik mengkonstruksi protokol SOAP 1.1 yang sangat tertib, disesuaikan persis dengan standar `.kredensial`:
- Menambahkan **`<iae:TeamID>`** dengan nilai `TEAM-25`.
- Mendeklarasikan aktivitas melalui **`<iae:ActivityName>`** dengan nilai `BookingCreated`.
- Melindungi dan membungkus *payload JSON* di dalam konstruksi **`<![CDATA[ ... ]]>`** (agar struktur JSON seperti kurung kurawal `{}` tidak dianggap sebagai tag XML).

## 2. Pemanfaatan Token M2M
Klien SOAP memanggil `GetM2MToken(ctx)` dari *helper* yang sebelumnya telah kita bangun. Hal ini memastikan setiap log audit dilindungi oleh token Bearer, menjadikan pertukaran data B2B (*Business-to-Business*) tetap aman.

## 3. Ekstraksi Nomor Resi (Receipt Number)
Respons dari peladen pusat dalam format XML diterima dan dilakukan proses *Unmarshaling* langsung ke struktur struct Go (`AuditResponseEnvelope`) dengan akurasi yang tinggi. Hal ini menjamin nilai resi seperti `IAE-LOG-2026-XXXXX` dapat dengan tepat diekstrak, dan sistem otomatis memberikan notifikasi kegagalan apabila status respons bukan `SUCCESS` atau format XML berubah.


# Prompt 4
Buat file rabbitmq_client.go di internal/infrastructure/. Buat fungsi untuk melakukan HTTP POST ke /api/v1/messages/publish. Gunakan M2M Token di header. Fungsi ini menerima struct event (misal: booking_id, status: LOCKED, timestamp), merubahnya menjadi JSON, dan menembaknya ke API gateway server cloud tersebut untuk mempublikasikan event ke iae.central.exchange secara asynchronous.