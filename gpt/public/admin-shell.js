import { createApp, reactive } from 'https://unpkg.com/vue@3/dist/vue.esm-browser.prod.js';

function createJsonFetcher() {
  return async (url, options = {}) => {
    const response = await fetch(url, {
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        ...(options.headers ?? {}),
      },
      ...options,
    });

    const contentType = response.headers.get('content-type') ?? '';
    const body = contentType.includes('application/json')
      ? await response.json()
      : await response.text();

    if (!response.ok) {
      throw new Error(
        typeof body === 'string'
          ? body
          : JSON.stringify(body),
      );
    }

    return body;
  };
}

const request = createJsonFetcher();

createApp({
  setup() {
    const mode = document.body.dataset.mode;
    const state = reactive({
      loading: false,
      error: '',
      message: '',
      admin: null,
      admins: [],
      users: [],
      clinics: [],
      algorithms: [],
      referenceData: { reasons: [], educationLevels: [] },
      login: { email: '', password: '' },
      register: {
        firstname: '',
        lastname: '',
        phone: '',
        email: '',
        verified: 'Y',
        userType: 'Super',
        password: '',
      },
      clinic: {
        name: '',
        address: '',
        latitude: '',
        longitude: '',
        added_by_id: '1',
      },
    });

    async function loadPanelData() {
      state.loading = true;
      state.error = '';
      try {
        const [admin, admins, users, clinics, reasons, educationLevels, algorithms] =
          await Promise.all([
            request('/getAdminDet'),
            request('/allAdmins'),
            request('/getUsers'),
            request('/getClinics'),
            request('/getContraceptionReason'),
            request('/getEducationalLevels'),
            request('/algo'),
          ]);
        state.admin = admin;
        state.admins = admins.data ?? [];
        state.users = users.data ?? [];
        state.clinics = clinics.data ?? [];
        state.algorithms = algorithms;
        state.referenceData.reasons = reasons;
        state.referenceData.educationLevels = educationLevels;
      } catch (error) {
        state.error = error instanceof Error ? error.message : String(error);
      } finally {
        state.loading = false;
      }
    }

    async function submitLogin() {
      state.loading = true;
      state.error = '';
      try {
        await request('/login', {
          method: 'POST',
          body: JSON.stringify(state.login),
        });
        window.location.assign('/panel');
      } catch (error) {
        state.error = error instanceof Error ? error.message : String(error);
      } finally {
        state.loading = false;
      }
    }

    async function submitRegister() {
      state.loading = true;
      state.error = '';
      try {
        await request('/admin', {
          method: 'POST',
          body: JSON.stringify(state.register),
        });
        window.location.assign('/panel');
      } catch (error) {
        state.error = error instanceof Error ? error.message : String(error);
      } finally {
        state.loading = false;
      }
    }

    async function createClinic() {
      state.loading = true;
      state.error = '';
      state.message = '';
      try {
        await request('/addClinic', {
          method: 'POST',
          body: JSON.stringify({
            ...state.clinic,
            latitude: Number(state.clinic.latitude),
            longitude: Number(state.clinic.longitude),
            added_by_id: Number(state.clinic.added_by_id),
          }),
        });
        state.message = 'Clinic created successfully.';
        state.clinic = {
          name: '',
          address: '',
          latitude: '',
          longitude: '',
          added_by_id: state.admin ? String(state.admin.id) : '1',
        };
        await loadPanelData();
      } catch (error) {
        state.error = error instanceof Error ? error.message : String(error);
      } finally {
        state.loading = false;
      }
    }

    if (mode === 'panel') {
      void loadPanelData();
    }

    return {
      mode,
      state,
      submitLogin,
      submitRegister,
      loadPanelData,
      createClinic,
    };
  },
  template: `
    <div class="shell">
      <header class="hero">
        <div>
          <p class="eyebrow">InCharge</p>
          <h1 v-if="mode === 'login'" class="icon-line"><i class="fa-solid fa-user-shield"></i><span>AdminloginComponent</span></h1>
          <h1 v-else-if="mode === 'register'" class="icon-line"><i class="fa-solid fa-shield-heart"></i><span>RegistersuperadminComponent</span></h1>
          <h1 v-else class="icon-line"><i class="material-icons">dashboard</i><span>AdminComponent</span></h1>
          <p class="subtitle">
            <span v-if="mode === 'login'">Verified admins sign in here.</span>
            <span v-else-if="mode === 'register'">Create the first Super admin to unlock the panel.</span>
            <span v-else>Session-backed admin panel for users, clinics, algorithms, and reference data.</span>
          </p>
        </div>
        <a v-if="mode === 'panel'" class="ghost" href="/logout"><i class="material-icons">logout</i><span>Logout</span></a>
      </header>

      <p v-if="state.error" class="status error">{{ state.error }}</p>
      <p v-if="state.message" class="status success">{{ state.message }}</p>

      <section v-if="mode === 'login'" class="card form-card">
        <label>
          <span>Email</span>
          <input v-model="state.login.email" type="email" autocomplete="username" />
        </label>
        <label>
          <span>Password</span>
          <input v-model="state.login.password" type="password" autocomplete="current-password" />
        </label>
        <button :disabled="state.loading" @click="submitLogin"><i class="material-icons">login</i><span>Sign In</span></button>
      </section>

      <section v-else-if="mode === 'register'" class="card form-card">
        <label><span>First name</span><input v-model="state.register.firstname" /></label>
        <label><span>Last name</span><input v-model="state.register.lastname" /></label>
        <label><span>Phone</span><input v-model="state.register.phone" /></label>
        <label><span>Email</span><input v-model="state.register.email" type="email" /></label>
        <label><span>Password</span><input v-model="state.register.password" type="password" /></label>
        <button :disabled="state.loading" @click="submitRegister"><i class="fa-solid fa-user-plus"></i><span>Create Super Admin</span></button>
      </section>

      <section v-else class="panel-grid">
        <article class="card overview">
          <h2>Current Admin</h2>
          <p v-if="state.admin">{{ state.admin.firstname }} {{ state.admin.lastname }} ({{ state.admin.userType }})</p>
          <p v-else-if="state.loading">Loading...</p>
          <button :disabled="state.loading" @click="loadPanelData"><i class="material-icons">refresh</i><span>Refresh Panel Data</span></button>
        </article>

        <article class="card form-card">
          <h2>Create Clinic</h2>
          <label><span>Name</span><input v-model="state.clinic.name" /></label>
          <label><span>Address</span><textarea v-model="state.clinic.address"></textarea></label>
          <label><span>Latitude</span><input v-model="state.clinic.latitude" type="number" step="0.0000001" /></label>
          <label><span>Longitude</span><input v-model="state.clinic.longitude" type="number" step="0.0000001" /></label>
          <label><span>Added by ID</span><input v-model="state.clinic.added_by_id" type="number" /></label>
          <button :disabled="state.loading" @click="createClinic"><i class="fa-solid fa-clinic-medical"></i><span>Create Clinic</span></button>
        </article>

        <article class="card">
          <h2>Admins</h2>
          <ul class="list">
            <li v-for="admin in state.admins" :key="admin.id">{{ admin.firstname }} {{ admin.lastname }} · {{ admin.verified }}</li>
          </ul>
        </article>

        <article class="card">
          <h2>Users</h2>
          <ul class="list">
            <li v-for="user in state.users" :key="user.id">{{ user.name }} · {{ user.email }}</li>
          </ul>
        </article>

        <article class="card">
          <h2>Clinics</h2>
          <ul class="list">
            <li v-for="clinic in state.clinics" :key="clinic.id">{{ clinic.name }} · {{ clinic.address }}</li>
          </ul>
        </article>

        <article class="card">
          <h2>Algorithms</h2>
          <ul class="list">
            <li v-for="algorithm in state.algorithms" :key="algorithm.id">#{{ algorithm.id }} · {{ algorithm.text }}</li>
          </ul>
        </article>

        <article class="card">
          <h2>Reference Data</h2>
          <p>Reasons: {{ state.referenceData.reasons.length }}</p>
          <p>Education Levels: {{ state.referenceData.educationLevels.length }}</p>
        </article>
      </section>
    </div>
  `,
}).mount('#app');
