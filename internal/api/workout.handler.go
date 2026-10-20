package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AryanBhatDev/Workout-App/internal/store"
	"github.com/AryanBhatDev/Workout-App/internal/utils"
	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *slog.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *slog.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {

	workoutIDInt, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workoutid"})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutByID(workoutIDInt)

	if err != nil {
		wh.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workoutid"})
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {

	var workout store.Workout

	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		wh.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workoutid"})
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)

	if err != nil {
		wh.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workoutid"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(createdWorkout)

}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")

	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}

	workoutIDInt, err := strconv.ParseInt(paramsWorkoutID, 10, 64)

	if err != nil {
		fmt.Println("e:", err.Error())
		http.NotFound(w, r)
		return
	}
	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutIDInt)

	if err != nil {
		http.Error(w, "failed to fetch workout", http.StatusInternalServerError)
		return
	}

	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)

	if err != nil {
		http.Error(w, "error while updating workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(existingWorkout)
}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")

	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}

	workoutIDInt, err := strconv.ParseInt(paramsWorkoutID, 10, 64)

	if err != nil {
		fmt.Println("e:", err.Error())
		http.NotFound(w, r)
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutIDInt)

	if err != nil {
		http.Error(w, "failed to fetch workout", http.StatusInternalServerError)
		return
	}

	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutIDInt)

	if err != nil {
		http.Error(w, "error while deleting workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNoContent)

}
