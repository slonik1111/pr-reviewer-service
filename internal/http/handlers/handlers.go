package handlers

import "net/http"

func RegisterAllHandlers(user *UserHandler, team *TeamHandler, pr *PRHandler) {
	// Users
	http.HandleFunc("/users/setIsActive", user.SetIsActive)
	http.HandleFunc("/users/getReview", user.GetPRsForUser)

	// Teams
	http.HandleFunc("/team/add", team.CreateTeam)
	http.HandleFunc("/team/get", team.GetTeam)

	// Pull Requests
	http.HandleFunc("/pullRequest/create", pr.CreatePR)
	http.HandleFunc("/pullRequest/merge", pr.MergePR)
	http.HandleFunc("/pullRequest/reassign", pr.ReassignReviewer)
}
