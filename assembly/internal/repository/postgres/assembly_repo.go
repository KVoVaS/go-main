package postgres

import (
	"context"
	"database/sql"
	"errors"

	"assembly/internal/domain"
	postgresgen "assembly/internal/repository/postgres/gen"
)

type AssemblyRepo struct {
	db *sql.DB
	q  *postgresgen.Queries
}

func NewAssemblyRepo(db *sql.DB) *AssemblyRepo {
	return &AssemblyRepo{db: db, q: postgresgen.New(db)}
}

func (r *AssemblyRepo) CreateEngine(ctx context.Context, e domain.Engine) (domain.Engine, error) {
	row, err := r.q.CreateEngine(ctx, postgresgen.CreateEngineParams{
		ID:         e.ID,
		Horsepower: e.Horsepower,
	})
	if err != nil {
		return domain.Engine{}, err
	}
	return domain.Engine{ID: row.ID, Horsepower: row.Horsepower}, nil
}

func (r *AssemblyRepo) CreateTransmission(ctx context.Context, t domain.Transmission) (domain.Transmission, error) {
	row, err := r.q.CreateTransmission(ctx, postgresgen.CreateTransmissionParams{
		ID:   t.ID,
		Type: t.Type,
	})
	if err != nil {
		return domain.Transmission{}, err
	}
	return domain.Transmission{ID: row.ID, Type: row.Type}, nil
}

func (r *AssemblyRepo) AssembleCar(ctx context.Context, car domain.Car, engineID, transID string) (domain.CarSpec, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.CarSpec{}, err
	}
	defer tx.Rollback()

	q := r.q.WithTx(tx)

	// Проверяем существование двигателя
	engine, err := q.GetEngine(ctx, engineID)
	if err != nil {
		return domain.CarSpec{}, errors.New("engine not found")
	}
	// Проверяем существование КПП
	trans, err := q.GetTransmission(ctx, transID)
	if err != nil {
		return domain.CarSpec{}, errors.New("transmission not found")
	}

	// Создаём автомобиль, если ещё не существует
	err = q.CreateCar(ctx, postgresgen.CreateCarParams{
		Vin:   car.VIN,
		Brand: car.Brand,
		Year:  car.Year,
	})
	if err != nil {
		return domain.CarSpec{}, err
	}

	// Связываем компоненты (поля EngineID и TransmissionID теперь sql.NullString)
	err = q.LinkComponents(ctx, postgresgen.LinkComponentsParams{
		Vin:            car.VIN,
		EngineID:       sql.NullString{String: engineID, Valid: true},
		TransmissionID: sql.NullString{String: transID, Valid: true},
	})
	if err != nil {
		return domain.CarSpec{}, err
	}

	if err = tx.Commit(); err != nil {
		return domain.CarSpec{}, err
	}

	return domain.CarSpec{
		Car:          car,
		Engine:       domain.Engine{ID: engine.ID, Horsepower: engine.Horsepower},
		Transmission: domain.Transmission{ID: trans.ID, Type: trans.Type},
	}, nil
}

func (r *AssemblyRepo) GetCarSpec(ctx context.Context, vin string) (domain.CarSpec, error) {
	row, err := r.q.GetCarSpec(ctx, vin)
	if err != nil {
		return domain.CarSpec{}, err
	}

	spec := domain.CarSpec{
		Car: domain.Car{
			VIN:   row.Vin,
			Brand: row.Brand,
			Year:  row.Year,
		},
	}
	if row.EngineID.Valid {
		spec.Engine.ID = row.EngineID.String
	}
	if row.Horsepower.Valid {
		spec.Engine.Horsepower = row.Horsepower.Int32
	}
	if row.TransID.Valid {
		spec.Transmission.ID = row.TransID.String
	}
	if row.TransType.Valid {
		spec.Transmission.Type = row.TransType.String
	}
	return spec, nil
}
