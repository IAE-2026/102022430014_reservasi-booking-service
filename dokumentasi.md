# Dokumentasi Penyesuaian Mesin Grader (Golang/Gin)

Dokumen ini berisi informasi yang diperlukan agar mesin **grader** dapat melakukan *static analysis*, validasi endpoint, dan pengujian otomatis terhadap aplikasi **Layanan Reservasi** yang dibangun menggunakan **Golang (Gin Framework)**.

## 1. Informasi Umum

* Framework: **Gin**
* Base URL: `http://localhost:7070`
* Prefix REST API: `/api/v1`
* Dokumentasi API: Swagger (Swaggo)
* GraphQL: gqlgen

---

## 2. Penyesuaian Static Analysis

Grader **tidak dapat menggunakan asumsi struktur Laravel** (`routes/api.php`), karena seluruh routing didaftarkan pada struktur proyek Go.

Lokasi file:

| Fungsi               | Lokasi                                      |
| -------------------- | ------------------------------------------- |
| Entry Point & Router | `cmd/app/main.go`                           |
| REST Endpoint        | `internal/delivery/rest/booking_handler.go` |
| GraphQL              | `internal/delivery/graphql/`                |

Jika melakukan pencarian route menggunakan parser atau regex, gunakan pola berikut:

| Endpoint   | Pola Gin                |
| ---------- | ----------------------- |
| Collection | `r.GET("/bookings"`     |
| Resource   | `r.GET("/bookings/:id"` |
| Action     | `r.POST("/bookings"`    |

> **Catatan:** Gin menggunakan parameter `:id`, bukan `{id}` seperti Laravel.

---

## 3. Endpoint yang Diuji

Semua endpoint REST menggunakan prefix `/api/v1`.

| Method | Endpoint               | Status Berhasil |
| ------ | ---------------------- | --------------- |
| GET    | `/api/v1/bookings`     | 200             |
| GET    | `/api/v1/bookings/:id` | 200 / 404       |
| POST   | `/api/v1/bookings`     | 201             |

Endpoint pendukung:

| Endpoint              | Keterangan         |
| --------------------- | ------------------ |
| `/swagger/index.html` | Swagger UI         |
| `/swagger/doc.json`   | OpenAPI JSON       |
| `/graphql`            | GraphQL Playground |
| `/graphql/v1/summary` | GraphQL Endpoint   |

---

## 4. Autentikasi

Seluruh endpoint REST memerlukan header:

```http
X-IAE-KEY: <NIM_MAHASISWA>
```

Perilaku yang diharapkan:

* Tanpa header → **401 Unauthorized**
* Dengan header valid → endpoint dapat diproses

---

## 5. Format Response

Response menggunakan wrapper JSON.

**Success**

```json
{
  "status": "success",
  "message": "...",
  "data": {}
}
```

**Error**

```json
{
  "status": "error",
  "message": "...",
  "errors": null
}
```

---

## 6. Aturan URL (Gin)

Gin menerapkan **Strict URL Matching**.

Gunakan:

```
GET /api/v1/bookings
```

Jangan gunakan:

```
GET /api/v1/bookings/
```

karena akan menghasilkan:

```
307 Temporary Redirect
```

Apabila HTTP Client secara otomatis menambahkan trailing slash, aktifkan **Follow Redirects**.

---

## 7. GraphQL

Endpoint GraphQL:

```
POST /graphql/v1/summary
```

Contoh introspection sederhana:

```graphql
query {
  __schema {
    queryType {
      name
    }
  }
}
```

---

## 8. Kriteria Pengujian

Mesin grader diharapkan melakukan validasi berikut:

| Pengujian             | Target                                |
| --------------------- | ------------------------------------- |
| API tanpa `X-IAE-KEY` | 401                                   |
| GET Collection        | 200                                   |
| GET Resource          | 200 / 404                             |
| POST Resource         | 201                                   |
| Swagger UI            | Dapat diakses                         |
| Swagger JSON          | Dapat diakses                         |
| GraphQL Playground    | Dapat diakses                         |
| GraphQL Query         | Berhasil                              |
| Service Docker        | Berjalan pada `http://localhost:7070` |

---

## 9. Perbedaan dengan Laravel

Agar proses evaluasi berjalan benar, mesin grader perlu memperhatikan perbedaan berikut:

* Route **tidak berada** pada `routes/api.php`.
* Gin menggunakan parameter route `:id`, bukan `{id}`.
* Gin menerapkan *Strict URL Matching* (307 pada trailing slash).
* Dokumentasi API menggunakan **Swaggo** (`/swagger/index.html`).
* Routing didaftarkan melalui kode Go pada `cmd/app/main.go`.
