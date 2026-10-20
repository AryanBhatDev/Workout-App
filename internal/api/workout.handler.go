package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/AryanBhatDev/Workout-App/internal/store"
	"github.com/AryanBhatDev/Workout-App/internal/utils"
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
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to fetch workout"})
		return
	}

	if workout == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {

	var workout store.Workout

	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		wh.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)

	if err != nil {
		wh.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": createdWorkout})

}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutIDInt, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workoutid"})
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutIDInt)

	if err != nil {
		wh.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to fetch workout"})
		return
	}

	if existingWorkout == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
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
		wh.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
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
		wh.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error while updating workout"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutIDInt, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workoutid"})
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutIDInt)

	if err != nil {
		wh.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to fetch workout"})
		return
	}

	if existingWorkout == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutIDInt)

	if err != nil {
		wh.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error while deleting workout"})
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
