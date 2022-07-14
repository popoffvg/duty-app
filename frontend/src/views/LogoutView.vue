<template>
  <div class="container">
    logout
  </div>

  <vue-basic-alert :duration="300" :closeIn="3000" ref="alertLogout" />
</template>

<script>
import axios from 'axios'
import VueBasicAlert from 'vue-basic-alert'

export default {
  components: {
    VueBasicAlert,
  },
  data() {
    return {};
  },
  created() {
    console.log("logout start");
    // logout
     axios
         .post("auth/logout",{ withCredentials: true })
         .then(
             response => {
               console.log("logout finish", response);

               if (response.data.status !== null) {
                 let status = response.data.status;
                 this.$refs.alertLogout.showAlert("success", status, "Success");
               }

               this.$router.push('/');
             },
             error => {
               console.log('Error logout: ', error);
               this.$refs.alertLogout.showAlert("error", error, "Error");
             },
         );
  },
  methods: {},
};
</script>