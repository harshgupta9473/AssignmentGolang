package database

import (
	"database/sql"

	_ "github.com/lib/pq"

	"log"
)

// func CreateTempUserTable(DB *sql.DB) error {
// 	query := `create table if not exists tempusers(
//         id serial primary key,
//         name varchar(255) not null,
//         email varchar(255) unique not null,
//         address text,
//         user_type varchar(50) not null check (user_type in ('Admin', 'Applicant')),
//         password_hash varchar(255) not null,
//         profile_headline text,
//         created_at timestamp default current_timestamp,
//         updated_at timestamp default current_timestamp,
// 		otp varchar(6) not null check (length(otp) = 6),
//         verified boolean default false,
//     )`
// 	_, err := DB.Exec(query)
// 	return err
// }

func CreateJob_ApplicationTable(DB *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS job_applications (
    id SERIAL PRIMARY KEY,
    job_id INTEGER REFERENCES jobs(id) ON DELETE CASCADE,
    applicant_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    applied_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`
	_, err := DB.Exec(query)
	return err
}

func CreateUsersTable(DB *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        address TEXT,
        user_type VARCHAR(50) NOT NULL CHECK (user_type IN ('admin', 'applicant')),
        password_hash VARCHAR(255) NOT NULL,
        profile_headline TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`
	_, err := DB.Exec(query)
	return err
}

func CreateProfilesTable(DB *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS profiles (
        id SERIAL PRIMARY KEY,
        user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE,
        resume_file_address TEXT,
        skills TEXT,
        education TEXT,
        experience TEXT,
        name VARCHAR(255),
        email VARCHAR(255),
        phone VARCHAR(20)
    )`
	_, err := DB.Exec(query)
	return err
}

func CreateJobsTable(DB *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS jobs (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        posted_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        total_applications INTEGER DEFAULT 0,
        company_name VARCHAR(255) NOT NULL,
        posted_by_id INTEGER REFERENCES users(id) ON DELETE SET NULL
    )`
	_, err := DB.Exec(query)
	return err
}

func InitTable(db *sql.DB) {
	err := CreateUsersTable(db)
	if err != nil {
		log.Fatal(err)
	}
	err = CreateProfilesTable(db)
	if err != nil {
		log.Fatal(err)
	}
	err = CreateJobsTable(db)
	if err != nil {
		log.Fatal(err)
	}
    err=CreateJob_ApplicationTable(db)
    if err!=nil{
        log.Fatal(err)
    }

}
