package model

import "time"

type Employee struct {
	ID               int64            `json:"ID" db:"id"`
	NIP              string           `json:"NIP" db:"nip"`
	NIPAlt           *string          `json:"NIP_ALT" db:"nip_alt"`
	NamaKaryawan     string           `json:"NAMA_KARYAWAN" db:"name"`
	Jabatan          string           `json:"JABATAN" db:"position"`
	KantorCabang     KantorCabang     `json:"KANTOR_CABANG" db:"branch_office"`
	IsExcluded       int64            `json:"IS_EXCLUDED" db:"is_excluded"`
	IsPanitia        int64            `json:"IS_PANITIA" db:"is_panitia"`
	JenisKepegawaian JenisKepegawaian `json:"JENIS_KEPEGAWAIAN" db:"employment_type"`
	Meja             *string          `json:"MEJA" db:"table"`
	PresentAt        *time.Time       `json:"PRESENT_AT" db:"present_at"`
}

type JenisKepegawaian string

const (
	Organik JenisKepegawaian = "Organik"
	Tad     JenisKepegawaian = "TAD"
)

type KantorCabang string

const (
	KantorCabangBogor        KantorCabang = "Kantor Cabang Bogor"
	KantorCabangCikarang     KantorCabang = "Kantor Cabang Cikarang"
	KantorCabangDepok        KantorCabang = "Kantor Cabang Depok"
	KantorCabangJakartaTimur KantorCabang = "Kantor Cabang Jakarta Timur"
	KantorCabangKarawang     KantorCabang = "Kantor Cabang Karawang"
	KantorCabangPurwokerto   KantorCabang = "Kantor Cabang Purwokerto"
	KantorCabangTangerang    KantorCabang = "Kantor Cabang Tangerang"
	KantorKasBekasi          KantorCabang = "Kantor Kas Bekasi"
	KantorPusatManajemen     KantorCabang = "Kantor Pusat Manajemen"
	KantorPusatOperasional   KantorCabang = "Kantor Pusat Operasional"
)
