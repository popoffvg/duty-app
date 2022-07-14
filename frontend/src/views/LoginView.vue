<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-6 offset-sm-3">
        <div>
          <h2>Login in Duty App</h2>
          <form @submit.prevent="handleSubmitLogin">
            <div class="form-group">
              <label>Username</label>
              <input type="text" v-model="username" name="username" class="form-control" :class="{ 'is-invalid': submitted && !username }" />
              <div v-show="submitted && !username" class="invalid-feedback">Username is required</div>
            </div>
            <div class="form-group">
              <label htmlFor="password">Password</label>
              <input type="password" v-model="password" name="password" class="form-control" :class="{ 'is-invalid': submitted && !password }" />
              <div v-show="submitted && !password" class="invalid-feedback">Password is required</div>
            </div>
            <div class="form-group">
              <button class="btn btn-primary" :disabled="loading">Login</button>
              <div>
                <img v-show="loading" src="data:image/gif;base64,R0lGODlhEAAQAPIAAP///wAAAMLCwkJCQgAAAGJiYoKCgpKSkiH/C05FVFNDQVBFMi4wAwEAAAAh/hpDcmVhdGVkIHdpdGggYWpheGxvYWQuaW5mbwAh+QQJCgAAACwAAAAAEAAQAAADMwi63P4wyklrE2MIOggZnAdOmGYJRbExwroUmcG2LmDEwnHQLVsYOd2mBzkYDAdKa+dIAAAh+QQJCgAAACwAAAAAEAAQAAADNAi63P5OjCEgG4QMu7DmikRxQlFUYDEZIGBMRVsaqHwctXXf7WEYB4Ag1xjihkMZsiUkKhIAIfkECQoAAAAsAAAAABAAEAAAAzYIujIjK8pByJDMlFYvBoVjHA70GU7xSUJhmKtwHPAKzLO9HMaoKwJZ7Rf8AYPDDzKpZBqfvwQAIfkECQoAAAAsAAAAABAAEAAAAzMIumIlK8oyhpHsnFZfhYumCYUhDAQxRIdhHBGqRoKw0R8DYlJd8z0fMDgsGo/IpHI5TAAAIfkECQoAAAAsAAAAABAAEAAAAzIIunInK0rnZBTwGPNMgQwmdsNgXGJUlIWEuR5oWUIpz8pAEAMe6TwfwyYsGo/IpFKSAAAh+QQJCgAAACwAAAAAEAAQAAADMwi6IMKQORfjdOe82p4wGccc4CEuQradylesojEMBgsUc2G7sDX3lQGBMLAJibufbSlKAAAh+QQJCgAAACwAAAAAEAAQAAADMgi63P7wCRHZnFVdmgHu2nFwlWCI3WGc3TSWhUFGxTAUkGCbtgENBMJAEJsxgMLWzpEAACH5BAkKAAAALAAAAAAQABAAAAMyCLrc/jDKSatlQtScKdceCAjDII7HcQ4EMTCpyrCuUBjCYRgHVtqlAiB1YhiCnlsRkAAAOwAAAAAAAAAAAA==" />
              </div>
            </div>
            <div v-if="error" class="alert alert-danger">{{error}}</div>
          </form>
        </div>
      </div>
    </div>
    </div>

  <!--
  <vue-basic-alert :duration="300" :closeIn="3000" ref="alert" />
  -->
</template>

<script>
import axios from 'axios'

export default {
  data() {
    return {
      username: "",
      password: "",
      submitted: false,
      loading: false,
      error: "",
      info: ""
    };
  },
  methods: {
    handleSubmitLogin() {
      this.submitted = true;
      const { username, password } = this;

      // stop here if form is invalid
      if (!(username && password)) {
        return;
      }

      this.loading = true;
      axios.post("auth/login", {}, {
        auth: {
          username: username,
          password: password
        }
      }).then(
          response => {
            this.loading = false;
            this.submitted = false;
            this.error = "";

            this.username = "";
            this.password = "";

            console.log('Authenticated', response);
            this.$router.push('/teams');
          },
          error => {
            console.log('Error on Authentication', error);
            this.error = error;
            this.loading = false;
            this.submitted = false;
          });
    }
  }
};
</script>