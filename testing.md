# Testing Guide - Tugas 2 Integrasi Aplikasi Enterprise (IAE-T2)

## Informasi Umum

Setiap service wajib memenuhi ketentuan pada dokumen **Standard Integration Contract (IAE-T2)** agar dapat terintegrasi dengan ekosistem Enterprise Application.

---

# Checklist Pengujian

## 1. Security & Standard

### 1.1 Validasi API Key

**Tujuan:**
Memastikan seluruh endpoint REST hanya dapat diakses menggunakan header autentikasi yang valid.

### Request Tanpa API Key

```http
GET /api/v1/[resource]
```

**Expected Result**

* Status Code: `401 Unauthorized`
* Response:

```json
{
  "status": "error",
  "message": "API Key is required",
  "errors": null
}
```

### Request Dengan API Key

```http
GET /api/v1/[resource]
X-IAE-KEY: [NIM]
```

**Expected Result**

* Status Code: `200 OK`
* Response menggunakan wrapper standar.

---

## 2. REST API Functional Testing

### 2.1 Collection Endpoint

**Endpoint**

```http
GET /api/v1/[resource]
```

**Header**

```http
X-IAE-KEY: [NIM]
```

**Expected Result**

* Status Code: `200 OK`
* Mengembalikan daftar data.

**Response Format**

```json
{
  "status": "success",
  "message": "Data retrieved successfully",
  "data": [],
  "meta": {
    "service_name": "[SERVICE_NAME]",
    "api_version": "v1"
  }
}
```

---

### 2.2 Resource Endpoint

**Endpoint**

```http
GET /api/v1/[resource]/{id}
```

**Header**

```http
X-IAE-KEY: [NIM]
```

### Data Ditemukan

**Expected Result**

* Status Code: `200 OK`

```json
{
  "status": "success",
  "message": "Data retrieved successfully",
  "data": {},
  "meta": {
    "service_name": "[SERVICE_NAME]",
    "api_version": "v1"
  }
}
```

### Data Tidak Ditemukan

**Expected Result**

* Status Code: `404 Not Found`

```json
{
  "status": "error",
  "message": "Resource not found",
  "errors": null
}
```

---

### 2.3 Action Endpoint

**Endpoint**

```http
POST /api/v1/[resource]
```

**Header**

```http
Content-Type: application/json
X-IAE-KEY: [NIM]
```

**Body**

```json
{
  "name": "Sample Data"
}
```

**Expected Result**

* Status Code: `201 Created`

```json
{
  "status": "success",
  "message": "Data created successfully",
  "data": {},
  "meta": {
    "service_name": "[SERVICE_NAME]",
    "api_version": "v1"
  }
}
```

---

## 3. Response Wrapper Validation

### Success Response

Seluruh response sukses wajib mengikuti format:

```json
{
  "status": "success",
  "message": "Deskripsi sukses",
  "data": {},
  "meta": {
    "service_name": "[SERVICE_NAME]",
    "api_version": "v1"
  }
}
```

### Error Response

Seluruh response error wajib mengikuti format:

```json
{
  "status": "error",
  "message": "Deskripsi error",
  "errors": null
}
```

---

## 4. Swagger/OpenAPI Testing

### Swagger UI Accessibility

**Tujuan**

Memastikan dokumentasi API dapat diakses.

**Endpoint Contoh**

```http
/swagger
```

atau

```http
/api/documentation
```

**Expected Result**

* Swagger UI dapat dibuka melalui browser.
* Tidak muncul error 404.
* Seluruh endpoint REST ditampilkan.

---

### Swagger Specification Validation

**Expected Result**

Swagger menampilkan:

* GET /api/v1/[resource]
* GET /api/v1/[resource]/{id}
* POST /api/v1/[resource]

Serta mendokumentasikan:

* Header `X-IAE-KEY`
* Request Body
* Response Schema

---

## 5. GraphQL Testing

### 5.1 Endpoint Accessibility

**Endpoint Contoh**

```http
/graphql
```

**Expected Result**

* Endpoint dapat diakses.
* Tidak menghasilkan 404.

---

### 5.2 GraphQL Playground

**Expected Result**

GraphQL Playground atau GraphiQL dapat dibuka.

Contoh URL:

```http
/graphql-playground
```

atau

```http
/playground
```

---

### 5.3 Introspection Query

**Query**

```graphql
{
  __schema {
    queryType {
      name
    }
  }
}
```

**Expected Result**

```json
{
  "data": {
    "__schema": {
      "queryType": {
        "name": "Query"
      }
    }
  }
}
```

---

### 5.4 Query Resource

Contoh Query:

```graphql
query {
  products {
    id
    name
  }
}
```

**Expected Result**

```json
{
  "data": {
    "products": [
      {
        "id": 1,
        "name": "Product A"
      }
    ]
  }
}
```

GraphQL harus memungkinkan client memilih field yang dibutuhkan.

---

## 6. Docker Deployment Testing

### Service Availability

**Tujuan**

Memastikan service berjalan di dalam container Docker.

### Verifikasi

```bash
docker ps
```

**Expected Result**

Container service dalam status:

```bash
Up
```

### Endpoint Test

```bash
curl http://localhost:[PORT]
```

**Expected Result**

Service memberikan response valid.

---

# Hasil Evaluasi Terakhir

| Kategori                                 | Status | Catatan                     |
| ---------------------------------------- | ------ | --------------------------- |
| Endpoint menolak request tanpa X-IAE-KEY | ✅      | Status 401 (Berhasil ditolak) |
| Request dengan X-IAE-KEY berhasil        | ✅      | Berhasil (Middleware membolehkan tanpa Bearer) |
| GET Collection Endpoint                  | ✅      | Berhasil, wrapper standar sesuai IAE-T2 |
| GET Resource Endpoint                    | ✅      | Berhasil, wrapper standar sesuai IAE-T2 |
| POST Action Endpoint                     | ✅      | Berhasil, autentikasi sukses dengan X-IAE-KEY |
| Swagger UI dapat diakses                 | ✅      | Berhasil (Ditemukan) |
| Swagger sesuai endpoint REST             | ✅      | Spec/path terbaca |
| GraphQL Endpoint dapat diakses           | ✅      | Berhasil diakses |
| Introspection Query berhasil             | ✅      | Berhasil mendapatkan schema |
| Service berjalan di Docker               | ✅      | Berhasil berjalan di Docker (Port 7070) |

---

# Definition of Done (DoD)

Service dinyatakan lulus apabila:

* [x] Endpoint Collection berjalan.
* [x] Endpoint Resource berjalan.
* [x] Endpoint Action berjalan.
* [ ] API Key `X-IAE-KEY` divalidasi.
* [x] Seluruh response menggunakan wrapper standar.
* [x] Swagger UI dapat diakses.
* [x] Swagger menampilkan seluruh endpoint.
* [x] GraphQL endpoint aktif.
* [x] GraphQL Playground aktif.
* [x] Minimal 1 Query GraphQL berfungsi.
* [x] Service berjalan pada Docker Container.
