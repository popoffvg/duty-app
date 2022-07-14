<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-8 offset-sm-2">
        <div>

          <div class="table-header">
            <h2>Duty history</h2>
          </div>

          <div class="duties-table">
            <b-table :items="history" :fields="fields">
              <template #cell(index)="data">
                <span>{{ data.index + 1}}</span>
              </template>
              <template #cell(name)="data">
                <span>{{ data.value }}</span>
              </template>
              <template #cell(date)="data">
                <span>{{ data.value }}</span>
              </template>
            </b-table>
          </div>

        </div>
      </div>
    </div>
  </div>

  <vue-basic-alert :duration="300" :closeIn="3000" ref="alertHistory" />
</template>

<script>
import axios from 'axios'
import VueBasicAlert from 'vue-basic-alert'
import util from '../assets/js/main.js';

export default {
  components: {
    VueBasicAlert,
  },
  data() {
    return {
      teamId: this.$route.params.teamId,
      fields: [
        { key: "index", label: "#"},
        { key: "name", label: "Name"},
        { key: "date", label: "Date"},
      ],
      history: [
        //{ "id": 5, "name": "test - 1! (new)", "is_daily": true, "date": "2022-05-31T00:00:00Z" },
        //{ "id": 3, "name": "test - 2! (new)",  "is_daily": false, "date": "2022-05-30T00:00:00Z" },
      ],
    };
  },
  created() {
    // load duties history
     axios
         .get("api/teams/" + this.teamId + "/history",{ withCredentials: true })
         .then(
             response => {
               if (response.data.data !== null) {
                 this.history = response.data.data;

                 for (let i = 0; i < this.history.length; i++) {
                   if (this.history[i].is_daily) {
                     this.history[i].date = util.formatDateDaily(this.history[i].date)
                   } else {
                     this.history[i].date = util.formatDateWeekly(this.history[i].date)
                   }
                 }

               } else {
                 this.$refs.alertHistory.showAlert("info", "History is empty", "Info");
               }
             },
             error => {
               console.log('Error load history: ', error);

               if (error.response.status === 401) {
                 this.$router.push('/');
                 return;
               }

               this.$refs.alertHistory.showAlert("error", error, "Error");
             },
         );
  },
  methods: {},
};
</script>