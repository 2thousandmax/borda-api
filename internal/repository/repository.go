package repository

import (
	"borda/internal/domain"
	"borda/internal/repository/postgres"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	SaveUser(username, password, contact string) (int, error)
	GetUserByCredentials(username, password string) (*domain.User, error)
	GetUserById(id int) (*domain.User, error)
	UpdatePassword(userId int, newPassword string) error
	AssignRole(userId, roleId int) error
	GetUserRole(userId int) (*domain.Role, error)
	// Зачем этот метод?
	//IsUsernameExists(username string) error

	// SetSession(userId int, session domain.Session) error
}

type TeamRepository interface {
	SaveTeam(teamLeaderId int, teamName string) (int, error)
	GetTeamById(teamId int) (*domain.Team, error)
	GetTeamByToken(token string) (*domain.Team, error)
	AddMember(teamId, userId int) error
	// Нужен ли нам этоти методы ???
	//GetMembers(teamId int) ([]domain.User, error)
	//IsTeamNameExists(teamName string) error
	//IsTeamTokenExists(token string) error May be unnecessary
	//IsTeamFull(teamId int) error
}

type TaskRepository interface {
	SaveTask(task domain.Task) (int, error)
	GetTaskById(id int) (*domain.Task, error)
	GetTasks(domain.TaskFilter) ([]*domain.Task, error)
	UpdateTask(id int, update domain.TaskUpdate) error
	SolveTask(taskId, teamId int) error
	SaveTaskSubmission(value domain.SubmitTaskRequest, isCorrect bool) error
	GetTaskSubmissions(taskId, teamId int) ([]*domain.TaskSubmission, error)
	// Для чего этот метод?
	//ChekSolvedTask(taskId, teamId int) (bool, error)
}

type SettingsRepository interface {
	Get(key string) (value string, err error)
	Set(key string, value string) (settingId int, err error)
}

type Repository struct {
	db       *sqlx.DB
	Users    UserRepository
	Teams    TeamRepository
	Tasks    TaskRepository
	Settings SettingsRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db:       db,
		Users:    postgres.NewUserRepository(db),
		Teams:    postgres.NewTeamRepository(db),
		Tasks:    postgres.NewTaskRepository(db),
		Settings: postgres.NewSettingsRepository(db),
	}
}
