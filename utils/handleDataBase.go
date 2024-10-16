package utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/harshgupta9473/recruitmentManagement/config"
	"github.com/harshgupta9473/recruitmentManagement/models"
	"golang.org/x/crypto/bcrypt"
)

func FindUserByEmail(email string) (*models.User, error) {
	// Prepare the SQL query
	db := config.GetDB()
	query := `SELECT id, name, email, address, user_type, password_hash, profile_headline, 
                     created_at, updated_at, otp, verified 
              FROM users 
              WHERE email = $1`

	// Execute the query
	row := db.QueryRow(query, email)

	// Create a TempUser to hold the result
	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.UserType,
		&user.Password, &user.ProfileHeadline, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// Return nil if no user was found
			return nil, nil
		}
		// Log any other error
		log.Println("Error scanning temp user:", err)
		return nil, err
	}

	return &user, nil
}

func InsertIntoUser(user *models.User) error {

	// Prepare the SQL query
	query := `INSERT INTO tempusers (name, email, address, user_type,encrypted_password, profile_headline, 
                     created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, current_timestamp, current_timestamp)`

	// Execute the query
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	db := config.GetDB()
	_, err := db.Exec(query, user.Name, user.Email, user.Address, user.UserType, user.Password,
		user.ProfileHeadline)

	if err != nil {
		log.Println("Error inserting or updating temp user:", err)
		return err
	}
	return nil
}

// func UpdateVerifiedStatus(email string, verified bool) error {
// 	// Prepare the SQL query to update the verified status
// 	query := `UPDATE tempusers
//               SET verified = $1, updated_at = current_timestamp
//               WHERE email = $2`

// 	db := config.GetDB()
// 	// Execute the query with the provided parameters
// 	_, err := db.Exec(query, verified, email)
// 	if err != nil {
// 		log.Println("Error updating verified status for email", email, ":", err)
// 		return err
// 	}

// 	return nil
// }

func InsertIntoJobs(jobreq *models.JobLisitingReq) error {
	query := `
    INSERT INTO jobs (title, description,posted_on,total_applications, company_name, posted_by_id)
    VALUES ($1, $2, $3, $4,$5)
    RETURNING id`

	db := config.GetDB()
	err := db.QueryRow(query, jobreq.Title, jobreq.Description, time.Now().UTC(), jobreq.TotalApplications, jobreq.CompanyName, jobreq.PostedByID)
	if err != nil {
		return fmt.Errorf("error inserting job: %w", err.Err())
	}

	return nil
}

func FindJobByID(id uint) (*models.Job, error) {
	var job *models.Job
	query := `select * from jobs where id=$1`
	db := config.GetDB()
	result := db.QueryRow(query, id)
	if result.Err() != nil {
		if result.Err() == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching the job")
	}
	err := result.Scan(&job.ID, &job.Title, &job.Description, &job.PostedOn, &job.TotalApplications, &job.CompanyName, &job.PostedByID)
	if err != nil {
		return nil, fmt.Errorf("error in scanning the job from result")
	}
	return job, nil
}

func InsertJobApplication(jobID, applicantID uint) error {
	query := `
    INSERT INTO job_applications (job_id, applicant_id, applied_on)
    VALUES ($1, $2, $3)`
	db := config.GetDB()
	_, err := db.Exec(query, jobID, applicantID, time.Now())
	if err != nil {
		return fmt.Errorf("error inserting job application: %w", err)
	}
	return nil
}

func GetApplicantsForJob(jobID uint) ([]models.User, error) {
	query := `
    SELECT u.*
    FROM job_applications ja
    JOIN users u ON ja.applicant_id = u.id
    WHERE ja.job_id = $1`

	db := config.GetDB()
	rows, err := db.Query(query, jobID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving applicants for job: %w", err)
	}
	defer rows.Close()

	var applicants []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.UserType,
			&user.Password, &user.ProfileHeadline, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning applicant details: %w", err)
		}
		applicants = append(applicants, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows result: %w", err)
	}

	return applicants, nil
}

func FetchAllUSers() ([]models.User, error) {
	query := `select * from users`
	db := config.GetDB()
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.UserType,
			&user.Password, &user.ProfileHeadline, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUSerBYID(id uint) (*models.User, error) {
	query := `select * from users where id=$1`
	db := config.GetDB()
	result := db.QueryRow(query, id)
	if result.Err() != nil {
		if result.Err() == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching the job")
	}
	var user models.User
	if err := result.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.UserType,
		&user.Password, &user.ProfileHeadline, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return &user,nil
}

func GetAllJobs()([]models.Job,error){
	query:=`select * from jobs`
	db := config.GetDB()
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []models.Job

	for rows.Next(){
		var job models.Job
		if err:= rows.Scan(&job.ID, &job.Title, &job.Description, &job.PostedOn, &job.TotalApplications, &job.CompanyName, &job.PostedByID); err!=nil{
			return nil,err
		}
		jobs=append(jobs, job)
	}
	if rows.Err()!=nil{
		return nil,err
	}
	return jobs,nil
}

func ApplyForJob(jobID uint, applicantID uint)error{
	db:=config.GetDB()

	tx,err:=db.Begin()
	if err!=nil{
		return err
	}

	defer func(){
		if err != nil {
            tx.Rollback()
            log.Println("Transaction rolled back:", err)
        }
	}()

	_, err = tx.Exec("INSERT INTO job_applications (job_id, applicant_id) VALUES ($1, $2)", jobID, applicantID)
    if err != nil {
        return err // Return the error which will trigger rollback
    }

	_, err = tx.Exec("UPDATE jobs SET total_applications = total_applications + 1 WHERE id = $1", jobID)
    if err != nil {
        return err // Return the error which will trigger rollback
    }
	if err:=tx.Commit(); err!=nil{
		return err
	}
	return nil
}



// type Job struct {
// 	ID                uint   // Primary key
// 	Title             string `json:"title"`
// 	Description       string `json:"description"`
// 	PostedOn          time.Time
// 	TotalApplications int
// 	CompanyName       string `json:"companyName"`
// 	PostedByID        uint   // Foreign key referring to User ID
// }
