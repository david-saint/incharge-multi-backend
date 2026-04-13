package handler

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/incharge/server/internal/config"
	"github.com/incharge/server/internal/dto"
	"github.com/incharge/server/internal/middleware"
	"github.com/incharge/server/internal/model"
	"github.com/incharge/server/internal/repository"
	"github.com/incharge/server/internal/service"
	"github.com/incharge/server/internal/validator"
)

// AdminHandler handles admin web panel endpoints.
type AdminHandler struct {
	adminRepo   *repository.AdminRepo
	userRepo    *repository.UserRepo
	clinicRepo  *repository.ClinicRepo
	algoRepo    *repository.AlgorithmRepo
	refRepo     *repository.ReferenceRepo
	authService *service.AuthService
	cfg         *config.Config
	templates   map[string]*template.Template
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(
	adminRepo *repository.AdminRepo,
	userRepo *repository.UserRepo,
	clinicRepo *repository.ClinicRepo,
	algoRepo *repository.AlgorithmRepo,
	refRepo *repository.ReferenceRepo,
	authService *service.AuthService,
	cfg *config.Config,
) *AdminHandler {
	h := &AdminHandler{
		adminRepo:   adminRepo,
		userRepo:    userRepo,
		clinicRepo:  clinicRepo,
		algoRepo:    algoRepo,
		refRepo:     refRepo,
		authService: authService,
		cfg:         cfg,
		templates:   make(map[string]*template.Template),
	}
	h.loadTemplates()
	return h
}

func (h *AdminHandler) loadTemplates() {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Join(filepath.Dir(filename), "..", "..", "templates", "admin")

	for _, name := range []string{"login", "register", "panel", "privacy"} {
		path := filepath.Join(baseDir, name+".html")
		t, err := template.ParseFiles(path)
		if err != nil {
			slog.Warn("admin template not found, using fallback", "template", name, "error", err)
			t = template.Must(template.New(name).Parse(adminFallbackTemplate(name)))
		}
		h.templates[name] = t
	}
}

func adminFallbackTemplate(name string) string {
	switch name {
	case "login":
		return `<!DOCTYPE html><html><head><title>InCharge Admin - Login</title>
<link href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css" rel="stylesheet">
</head><body class="grey lighten-4"><div class="container"><div class="row" style="margin-top:100px"><div class="col s12 m6 offset-m3"><div class="card"><div class="card-content"><span class="card-title center-align">InCharge Admin Login</span>
<form id="loginForm" method="POST" action="/login"><div class="input-field"><input type="email" name="email" id="email" required><label for="email">Email</label></div>
<div class="input-field"><input type="password" name="password" id="password" required><label for="password">Password</label></div>
<button class="btn waves-effect waves-light blue darken-2 col s12" type="submit">Login</button></form>
<div id="error" class="red-text center-align" style="margin-top:10px"></div></div></div></div></div></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/js/materialize.min.js"></script>
<script>document.getElementById('loginForm').addEventListener('submit',function(e){e.preventDefault();
fetch('/login',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({email:document.getElementById('email').value,password:document.getElementById('password').value})})
.then(r=>{if(r.ok)window.location='/panel';else document.getElementById('error').textContent='Invalid credentials or unverified account';})
.catch(()=>document.getElementById('error').textContent='Login failed');});</script></body></html>`
	case "register":
		return `<!DOCTYPE html><html><head><title>InCharge Admin - Setup</title>
<link href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css" rel="stylesheet">
</head><body class="grey lighten-4"><div class="container"><div class="row" style="margin-top:60px"><div class="col s12 m8 offset-m2"><div class="card"><div class="card-content"><span class="card-title center-align">Create Super Admin</span>
<form id="regForm"><div class="row"><div class="input-field col s6"><input type="text" name="firstname" id="firstname" required><label for="firstname">First Name</label></div>
<div class="input-field col s6"><input type="text" name="lastname" id="lastname" required><label for="lastname">Last Name</label></div></div>
<div class="input-field"><input type="email" name="email" id="email" required><label for="email">Email</label></div>
<div class="input-field"><input type="text" name="phone" id="phone"><label for="phone">Phone</label></div>
<div class="input-field"><input type="password" name="password" id="password" required><label for="password">Password (min 6 chars)</label></div>
<button class="btn waves-effect waves-light blue darken-2 col s12" type="submit">Create Admin</button></form>
<div id="error" class="red-text center-align" style="margin-top:10px"></div></div></div></div></div></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/js/materialize.min.js"></script>
<script>document.getElementById('regForm').addEventListener('submit',function(e){e.preventDefault();
fetch('/admin',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({firstname:document.getElementById('firstname').value,lastname:document.getElementById('lastname').value,email:document.getElementById('email').value,phone:document.getElementById('phone').value,password:document.getElementById('password').value,verified:'Y',userType:'Super'})})
.then(r=>r.json()).then(d=>{if(d.id)window.location='/panel';else document.getElementById('error').textContent=d.message||'Registration failed';})
.catch(()=>document.getElementById('error').textContent='Registration failed');});</script></body></html>`
	case "panel":
		return `<!DOCTYPE html><html><head><title>InCharge Admin Panel</title>
<link href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css" rel="stylesheet">
<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
<link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css" rel="stylesheet">
<style>body{display:flex;min-height:100vh;flex-direction:column}main{flex:1 0 auto}.sidenav{width:250px}.sidenav .user-view{padding:20px}.content{margin-left:250px;padding:20px}@media only screen and (max-width:992px){.content{margin-left:0}}</style>
</head><body><div id="app">
<ul class="sidenav sidenav-fixed blue darken-3"><li><div class="user-view"><span class="white-text name">InCharge Admin</span></div></li>
<li><a class="white-text" href="#" onclick="loadSection('dashboard')"><i class="material-icons white-text">dashboard</i>Dashboard</a></li>
<li><a class="white-text" href="#" onclick="loadSection('users')"><i class="material-icons white-text">people</i>Users</a></li>
<li><a class="white-text" href="#" onclick="loadSection('clinics')"><i class="material-icons white-text">local_hospital</i>Clinics</a></li>
<li><a class="white-text" href="#" onclick="loadSection('admins')"><i class="material-icons white-text">admin_panel_settings</i>Admins</a></li>
<li><a class="white-text" href="#" onclick="loadSection('algorithms')"><i class="material-icons white-text">account_tree</i>Algorithms</a></li>
<li><div class="divider"></div></li>
<li><a class="white-text" href="/logout"><i class="material-icons white-text">exit_to_app</i>Logout</a></li></ul>
<main class="content"><div id="content"><h4>Welcome to InCharge Admin Panel</h4><p>Select a section from the sidebar to manage.</p></div></main></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/js/materialize.min.js"></script>
<script>
function loadSection(s){const c=document.getElementById('content');
switch(s){
case 'users':fetch('/getUsers').then(r=>r.json()).then(d=>renderTable(c,'Users',d.data||[],'id','name','email'));break;
case 'clinics':fetch('/getClinics').then(r=>r.json()).then(d=>renderTable(c,'Clinics',d.data||[],'id','name','address'));break;
case 'admins':fetch('/allAdmins').then(r=>r.json()).then(d=>renderTable(c,'Admins',d.data||[],'id','firstname','email','verified'));break;
case 'algorithms':fetch('/algo').then(r=>r.json()).then(d=>renderTable(c,'Algorithms',d||[],'id','text','actionType','active'));break;
default:c.innerHTML='<h4>Dashboard</h4><p>Welcome to the admin panel.</p>';}}
function renderTable(c,title,rows,...cols){let h='<h4>'+title+'</h4>';if(!rows.length){c.innerHTML=h+'<p>No records found.</p>';return;}
h+='<table class="striped highlight"><thead><tr>';cols.forEach(c=>h+='<th>'+c+'</th>');h+='</tr></thead><tbody>';
rows.forEach(r=>{h+='<tr>';cols.forEach(c=>h+='<td>'+(r[c]!=null?r[c]:'')+'</td>');h+='</tr>';});
h+='</tbody></table>';c.innerHTML=h;}
</script></body></html>`
	case "privacy":
		return `<!DOCTYPE html><html><head><title>InCharge - Privacy Policy</title>
<link href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css" rel="stylesheet">
</head><body><div class="container" style="margin-top:40px"><h3>Privacy Policy</h3>
<p>InCharge is committed to protecting your privacy. This policy describes how we collect, use, and share your personal information.</p>
<h5>Information We Collect</h5><p>We collect information you provide when creating an account, including your name, email, phone number, and health-related profile data for contraceptive guidance.</p>
<h5>How We Use Your Information</h5><p>Your information is used solely to provide personalized contraceptive guidance and connect you with nearby clinics. We do not sell your data to third parties.</p>
<h5>Data Security</h5><p>We use industry-standard encryption and security measures to protect your data.</p>
<h5>Contact</h5><p>For questions about this policy, please contact us at support@incharge.app.</p></div></body></html>`
	default:
		return `<html><body><p>Page not found</p></body></html>`
	}
}

// --- Web Routes ---

// Index handles GET / — registration or redirect.
func (h *AdminHandler) Index(w http.ResponseWriter, r *http.Request) {
	if !h.adminRepo.SuperAdminExists() {
		h.templates["register"].Execute(w, nil)
		return
	}
	http.Redirect(w, r, "/loginView", http.StatusFound)
}

// LoginView handles GET /loginView.
func (h *AdminHandler) LoginView(w http.ResponseWriter, r *http.Request) {
	h.templates["login"].Execute(w, nil)
}

// Login handles POST /login — admin login.
func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.AdminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	admin, err := h.adminRepo.FindByEmail(req.Email)
	if err != nil {
		dto.WriteJSON(w, 501, map[string]string{"message": "Invalid credentials"})
		return
	}

	if !h.authService.CheckPassword(admin.Password, req.Password) {
		dto.WriteJSON(w, 501, map[string]string{"message": "Invalid credentials"})
		return
	}

	if !admin.IsVerified() {
		dto.WriteJSON(w, 501, map[string]string{"message": "Account not verified"})
		return
	}

	if err := middleware.SetAdminSession(w, r, admin.ID); err != nil {
		dto.WriteServerError(w, "Session error", err, h.cfg.App.IsProduction)
		return
	}

	dto.WriteJSON(w, http.StatusOK, map[string]string{"message": "Login successful"})
}

// Logout handles GET /logout.
func (h *AdminHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := middleware.ClearAdminSession(w, r); err != nil {
		slog.Warn("failed to clear admin session", "error", err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// Panel handles GET /panel.
func (h *AdminHandler) Panel(w http.ResponseWriter, r *http.Request) {
	h.templates["panel"].Execute(w, nil)
}

// Privacy handles GET /privacy.
func (h *AdminHandler) Privacy(w http.ResponseWriter, r *http.Request) {
	h.templates["privacy"].Execute(w, nil)
}

// --- Admin Management ---

// CreateAdmin handles POST /admin.
// When no Super admin exists, this endpoint is accessible without a session (first-time setup).
// Otherwise, it requires a verified admin session.
func (h *AdminHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	// If a Super admin already exists, enforce session-based auth.
	if h.adminRepo.SuperAdminExists() {
		session, err := middleware.SessionStore.Get(r, "incharge_session")
		if err != nil || session.IsNew {
			dto.WriteAuthError(w, "Permission Denied")
			return
		}
		adminID, ok := session.Values["admin_id"].(uint)
		if !ok || adminID == 0 {
			dto.WriteAuthError(w, "Permission Denied")
			return
		}
		// Verify the session admin is verified.
		admin, err := h.adminRepo.FindByID(adminID)
		if err != nil || !admin.IsVerified() {
			dto.WriteAuthError(w, "Permission Denied")
			return
		}
	}

	var req dto.AdminCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	hash, err := h.authService.HashPassword(req.Password)
	if err != nil {
		dto.WriteServerError(w, "Failed to hash password", err, h.cfg.App.IsProduction)
		return
	}

	accessToken, err := service.GenerateRandomToken(32)
	if err != nil {
		dto.WriteServerError(w, "Failed to generate token", err, h.cfg.App.IsProduction)
		return
	}

	verified := req.Verified
	if verified == "" {
		verified = "N"
	}

	admin := &model.Admin{
		Firstname:   req.Firstname,
		Lastname:    req.Lastname,
		Email:       req.Email,
		Phone:       model.NewNullString(req.Phone, req.Phone != ""),
		Password:    hash,
		Verified:    verified,
		UserType:    req.UserType,
		AccessToken: model.NewNullString(accessToken, true),
	}

	if err := h.adminRepo.Create(admin); err != nil {
		dto.WriteServerError(w, "Failed to create admin", err, h.cfg.App.IsProduction)
		return
	}

	// If first Super admin, auto-login.
	if admin.UserType == "Super" && admin.IsVerified() {
		middleware.SetAdminSession(w, r, admin.ID)
	}

	dto.WriteJSON(w, http.StatusCreated, admin)
}

// ListAdmins handles GET /allAdmins.
func (h *AdminHandler) ListAdmins(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	admins, total, err := h.adminRepo.ListPaginated(page, 50)
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch admins", err, h.cfg.App.IsProduction)
		return
	}

	path := h.cfg.App.URL + r.URL.Path
	resp := dto.NewPaginatedResponse(admins, total, page, 50, path)
	dto.WriteJSON(w, http.StatusOK, resp)
}

// GetAdminDetails handles GET /getAdminDet.
func (h *AdminHandler) GetAdminDetails(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetAdminID(r.Context())
	if !ok {
		dto.WriteAuthError(w, "Permission Denied")
		return
	}

	admin, err := h.adminRepo.FindByID(adminID)
	if err != nil {
		dto.WriteAuthError(w, "Permission Denied")
		return
	}

	dto.WriteJSON(w, http.StatusOK, admin)
}

// UpdateAdmin handles PUT /updateAdmin/{admin_id}.
func (h *AdminHandler) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "admin_id"), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid admin ID"})
		return
	}

	admin, err := h.adminRepo.FindByID(uint(id))
	if err != nil {
		dto.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "Admin not found"})
		return
	}

	var req dto.AdminUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	if req.Firstname != nil {
		admin.Firstname = *req.Firstname
	}
	if req.Lastname != nil {
		admin.Lastname = *req.Lastname
	}
	if req.Phone != nil {
		admin.Phone = model.NewNullString(*req.Phone, *req.Phone != "")
	}
	if req.Email != nil {
		admin.Email = *req.Email
	}
	if req.Verified != nil {
		admin.Verified = *req.Verified
	}
	if req.UserType != nil {
		admin.UserType = *req.UserType
	}
	if req.AccessToken != nil {
		admin.AccessToken = model.NewNullString(*req.AccessToken, *req.AccessToken != "")
	}

	if err := h.adminRepo.Update(admin); err != nil {
		dto.WriteServerError(w, "Failed to update admin", err, h.cfg.App.IsProduction)
		return
	}

	dto.WriteJSON(w, http.StatusOK, admin)
}

// --- Admin panel user management ---

// AdminListUsers handles GET /getUsers.
func (h *AdminHandler) AdminListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	users, total, err := h.userRepo.ListWithProfile(page, 50)
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch users", err, h.cfg.App.IsProduction)
		return
	}

	resources := make([]dto.UserResource, len(users))
	for i, u := range users {
		resources[i] = toUserResource(&u)
	}

	path := h.cfg.App.URL + r.URL.Path
	resp := dto.NewPaginatedResponse(resources, total, page, 50, path)
	dto.WriteJSON(w, http.StatusOK, resp)
}

// AdminDeleteUser handles DELETE /deleteUser/{user_id}.
func (h *AdminHandler) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
		return
	}
	if err := h.userRepo.SoftDelete(uint(id)); err != nil {
		dto.WriteServerError(w, "Failed to delete user", err, h.cfg.App.IsProduction)
		return
	}
	dto.WriteSuccess(w, "User deleted successfully", nil)
}

// AdminRestoreUser handles PUT /updateUser/{user_id} and PUT /revertDeletedUser/{user_id}.
func (h *AdminHandler) AdminRestoreUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
		return
	}
	if err := h.userRepo.Restore(uint(id)); err != nil {
		dto.WriteServerError(w, "Failed to restore user", err, h.cfg.App.IsProduction)
		return
	}
	dto.WriteSuccess(w, "User restored successfully", nil)
}

// AdminListDeletedUsers handles GET /getDeletedUsers.
func (h *AdminHandler) AdminListDeletedUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.ListDeleted()
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch deleted users", err, h.cfg.App.IsProduction)
		return
	}

	resources := make([]dto.UserResource, len(users))
	for i, u := range users {
		resources[i] = toUserResource(&u)
	}
	dto.WriteJSON(w, http.StatusOK, resources)
}

// --- Admin panel clinic management ---

// AdminListClinics handles GET /getClinics.
func (h *AdminHandler) AdminListClinics(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	clinics, total, err := h.clinicRepo.ListPaginated(page, 50)
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch clinics", err, h.cfg.App.IsProduction)
		return
	}

	resources := make([]dto.ClinicResource, len(clinics))
	for i, c := range clinics {
		resources[i] = toClinicResource(&c)
	}

	path := h.cfg.App.URL + r.URL.Path
	resp := dto.NewPaginatedResponse(resources, total, page, 50, path)
	dto.WriteJSON(w, http.StatusOK, resp)
}

// AdminAddClinic handles POST /addClinic.
func (h *AdminHandler) AdminAddClinic(w http.ResponseWriter, r *http.Request) {
	h.AddClinicInternal(w, r)
}

// AdminUpdateClinic handles PUT /updateClinic/{clinic_id}.
func (h *AdminHandler) AdminUpdateClinic(w http.ResponseWriter, r *http.Request) {
	// Remap param name for admin routes.
	h.UpdateClinicInternal(w, r, "clinic_id")
}

// AdminDeleteClinic handles DELETE /deleteClinic/{clinic_id}.
func (h *AdminHandler) AdminDeleteClinic(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "clinic_id"), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid clinic ID"})
		return
	}
	if err := h.clinicRepo.SoftDelete(uint(id)); err != nil {
		dto.WriteServerError(w, "Failed to delete clinic", err, h.cfg.App.IsProduction)
		return
	}
	dto.WriteSuccess(w, "Clinic deleted successfully", nil)
}

// AdminListDeletedClinics handles GET /getDeletedClinics.
func (h *AdminHandler) AdminListDeletedClinics(w http.ResponseWriter, r *http.Request) {
	clinics, err := h.clinicRepo.ListDeleted()
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch deleted clinics", err, h.cfg.App.IsProduction)
		return
	}
	resources := make([]dto.ClinicResource, len(clinics))
	for i, c := range clinics {
		resources[i] = toClinicResource(&c)
	}
	dto.WriteJSON(w, http.StatusOK, resources)
}

// AdminRestoreClinic handles PUT /revertDeletedClinic/{clinic_id}.
func (h *AdminHandler) AdminRestoreClinic(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "clinic_id"), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid clinic ID"})
		return
	}
	if err := h.clinicRepo.Restore(uint(id)); err != nil {
		dto.WriteServerError(w, "Failed to restore clinic", err, h.cfg.App.IsProduction)
		return
	}
	dto.WriteSuccess(w, "Clinic restored successfully", nil)
}

// --- Shared internal helpers ---

func (h *AdminHandler) AddClinicInternal(w http.ResponseWriter, r *http.Request) {
	var req dto.ClinicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}
	if errs := validator.ValidateStruct(req); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}
	clinic := &model.Clinic{
		Name:      req.Name,
		Address:   req.Address,
		Latitude:  &req.Latitude,
		Longitude: &req.Longitude,
		AddedByID: req.AddedByID,
	}
	if err := h.clinicRepo.Create(clinic); err != nil {
		dto.WriteServerError(w, "Failed to create clinic", err, h.cfg.App.IsProduction)
		return
	}
	dto.WriteJSON(w, http.StatusCreated, dto.SuccessResponse{
		Status: true, Message: "Clinic created successfully", Data: toClinicResource(clinic),
	})
}

func (h *AdminHandler) UpdateClinicInternal(w http.ResponseWriter, r *http.Request, paramName string) {
	id, err := strconv.ParseUint(chi.URLParam(r, paramName), 10, 64)
	if err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid clinic ID"})
		return
	}
	clinic, err := h.clinicRepo.FindByID(uint(id))
	if err != nil {
		dto.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "Clinic not found"})
		return
	}
	var req dto.ClinicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}
	clinic.Name = req.Name
	clinic.Address = req.Address
	clinic.Latitude = &req.Latitude
	clinic.Longitude = &req.Longitude
	if err := h.clinicRepo.Update(clinic); err != nil {
		dto.WriteServerError(w, "Failed to update clinic", err, h.cfg.App.IsProduction)
		return
	}
	dto.WriteSuccess(w, "Clinic updated successfully", toClinicResource(clinic))
}

// --- Reference data for admin ---

// AdminListContraceptionReasons handles GET /getContraceptionReason.
func (h *AdminHandler) AdminListContraceptionReasons(w http.ResponseWriter, r *http.Request) {
	reasons, err := h.refRepo.ListContraceptionReasons()
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch reasons", err, h.cfg.App.IsProduction)
		return
	}
	resources := make([]dto.ReasonResource, len(reasons))
	for i, reason := range reasons {
		resources[i] = dto.ReasonResource{ID: reason.ID, Value: reason.Value}
	}
	dto.WriteJSON(w, http.StatusOK, resources)
}

// AdminListEducationLevels handles GET /getEducationalLevels.
func (h *AdminHandler) AdminListEducationLevels(w http.ResponseWriter, r *http.Request) {
	levels, err := h.refRepo.ListEducationLevels()
	if err != nil {
		dto.WriteServerError(w, "Failed to fetch education levels", err, h.cfg.App.IsProduction)
		return
	}
	resources := make([]dto.NamedResource, len(levels))
	for i, level := range levels {
		resources[i] = dto.NamedResource{ID: level.ID, Name: level.Name}
	}
	dto.WriteJSON(w, http.StatusOK, resources)
}
