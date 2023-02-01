<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-8 offset-sm-2">

        <div>
          <b-collapse v-model="showDuties">
            <div class="table-header">
            <h2>Duties</h2>
          </div>

            <div class="duties-container">
              <b-row class="text-center">
                <b-col class="duty-type-col text-left">Daily</b-col>
                <b-col class="duty-date-col">{{ dailyDate }}</b-col>
                <b-col class="duty-teammate-col">
                  <b-form-select v-on:change="changeDutyTeammate($event, true)" v-model="selectedDailyTeammateId" value-field="id" text-field="name" :options="teammates" class="form-control"></b-form-select>
                </b-col>
                <b-col class="duty-delete-col">
                  <b-button class="delete-duty" @click="deleteDuty(true)">
                    <span>X</span>
                  </b-button>
                </b-col>
              </b-row>
              <div class="w-100 duty-spacer"></div>
              <b-row class="text-center">
                <b-col class="duty-type-col text-left">Weekly</b-col>
                <b-col class="duty-date-col">{{ weeklyDate }}</b-col>
                <b-col class="duty-teammate-col">
                  <b-form-select v-on:change="changeDutyTeammate($event, false)" v-model="selectedWeeklyTeammateId" value-field="id" text-field="name" :options="teammates" class="form-control"></b-form-select>
                </b-col>
                <b-col class="duty-delete-col">
                  <b-button class="delete-duty" @click="deleteDuty(false)">
                    <span>X</span>
                  </b-button>
                </b-col>
              </b-row>


              <div class="history-block">
                <a :href="'#/teams/' + teamId + '/settings'">settings</a>
                &nbsp;
                <a :href="'#/teams/' + teamId + '/history'">history</a>
              </div>
          </div>
          </b-collapse>
        </div>


        <div>
          <div class="table-header">
            <h2>Teammates</h2>
          </div>

          <div class="add-button">
            <b-button @click="addTeammate()">
              <span>Add teammate</span>
            </b-button>
          </div>

          <div class="teammates-table">
            <b-table :items="teammates" :fields="fields">
              <template #cell(index)="data">
                <span>{{ data.index + 1}}</span>
              </template>
              <template #cell(name)="data">
                <b-form-input v-if="teammates[data.index].isEdit" type="text" v-model="teammates[data.index].name"></b-form-input>
                <span v-else>{{ data.value }}</span>
              </template>
              <template #cell(duty_readiness)="data">
                <b-form-checkbox :disabled="!teammates[data.index].isEdit" v-model="teammates[data.index].duty_readiness"></b-form-checkbox>
              </template>
              <template #cell(duties)="data">
                <b-form-input v-if="teammates[data.index].isEdit" type="text" v-model="teammates[data.index].duties"></b-form-input>
                <span v-else>{{ data.value }}</span>
              </template>
              <template #cell(options)="data">
                <div class="table-buttons">
                  <b-button @click="editRowHandler(data)">
                    <span v-if="!teammates[data.index].isEdit">Edit</span>
                    <span v-else>Save</span>
                  </b-button>
                  &nbsp;
                  <b-button @click="deleteTeammateRowHandler(data)" class="delete-button">
                    <span>Delete</span>
                  </b-button>
                </div>
              </template>
            </b-table>
          </div>

        </div>
      </div>
    </div>
  </div>

  <div v-if="showModalTeammate" class="modal">
    <!-- Modal content -->
    <div class="modal-content">
      <span class="close" @click="cancelDeleteTeammate()">&times;</span>
      <p>Are you sure to delete teammate '{{ dataToDelete.name }}'?</p>

      <div class="modal-buttons">
        <b-button @click="cancelDeleteTeammate()">
          <span>cancel</span>
        </b-button>
        <b-button class="delete-button" @click="deleteTeammate()">
          <span>delete</span>
        </b-button>
      </div>
    </div>
  </div>

  <vue-basic-alert :duration="300" :closeIn="3000" ref="alertTeammate" />
</template>

<style>
  .teammates-table thead, tbody, tfoot, tr, td, th {
    text-align: left;
    vertical-align: middle !important;
  }
  .teammates-table th:nth-of-type(3) {
    text-align: center;
  }
  .teammates-table th:nth-of-type(4) {
    text-align: center;
  }
  .teammates-table th:nth-of-type(5) {
    text-align: center;
  }
  .teammates-table td:nth-of-type(1) {
    width:40px;
  }
  .teammates-table td:nth-of-type(3) {
    width:80px;
    text-align: center;
  }
  .teammates-table td:nth-of-type(4) {
    width:80px;
    text-align: center;
  }
  .teammates-table td:nth-of-type(5) {
    width:165px;
  }
</style>

<script>
import axios from 'axios'
import util from '../assets/js/main.js';


export default {
  data() {
    return {
      dailyDate: "today",
      weeklyDate: "this week",
      dailyDutyId: null,
      weeklyDutyId: null,
      selectedDailyTeammateId: null,
      selectedWeeklyTeammateId: null,
      showDuties: true,
      duties: [
        //{ "id": 5, "team_id": 18, "teammate_id": 18, "is_daily": false, "date": "2022-05-16T00:00:00Z" },
        //{ "id": 6, "team_id": 18, "teammate_id": 19, "is_daily": true, "date": "2022-05-22T00:00:00Z" }
      ],
      showModalTeammate: false,
      dataToDelete: {},
      teamId: this.$route.params.teamId,
      fields: [
        { key: "index", label: "#"},
        { key: "name", label: "Name"},
        { key: "duty_readiness", label: "Ready?"},
        { key: "duties", label: "W.Duties"},
        { key: "options", label: "Options"},
      ],
      teammates: [
       //{ id: 18, team_id: 2, name: "Teammate 1", duty_readiness: "true", duties: 2 },
       //{ id: 19, team_id: 2, name: "Teammate 2", duty_readiness: "false", duties: 3 },
      ],
    };
  },
  created() {
    // list all teammates
     axios
         .get("api/teams/" + this.teamId + "/teammates",{ withCredentials: true })
         .then(
             response => {
               if (response.data.data !== null) {
                 this.teammates = response.data.data;
                 this.teammates = this.teammates.map(item => ({...item, isEdit: false}));
                 this.showDuties = true;
               } else {
                 this.$refs.alertTeammate.showAlert("info", "Click on [Add teammate] button", "Info");
               }
             },
             error => {
               console.log('Error load teams: ', error);

               if (error.response.status === 401) {
                 this.$router.push('/');
                 return;
               }

               this.$refs.alertTeammate.showAlert("error", error, "Error");
             },
         );

    // get daily and weekly duties
    axios
        .get("api/teams/" + this.teamId + "/duties/now",{ withCredentials: true })
        .then(
            response => {
              if (response.data.data !== null) {
                this.duties = response.data.data;
                this.duties.forEach(duty => {
                  if (duty.is_daily) {
                    this.selectedDailyTeammateId = duty.teammate_id;
                    this.dailyDutyId = duty.id;

                    console.log("daily", util.formatDateDaily(duty.date));
                    this.dailyDate = util.formatDateDaily(duty.date);
                  } else {
                    this.selectedWeeklyTeammateId = duty.teammate_id;
                    this.weeklyDutyId = duty.id;

                    console.log("weekly", util.formatDateWeekly(duty.date));
                    this.weeklyDate = util.formatDateWeekly(duty.date);
                  }
                });
              } else {
                if (this.teammates != null) {
                  this.$refs.alertTeammate.showAlert("info", "No duties today", "Info");
                }
              }
            },
            error => {
              console.log('Error load duties: ', error);
              this.$refs.alertTeammate.showAlert("error", error, "Error");
            },
        );
  },
  methods: {
    editRowHandler(data) {
      this.teammates[data.index].isEdit = !this.teammates[data.index].isEdit;

      if (!this.teammates[data.index].isEdit) {
        // validate duties is number
        let duties = parseInt(this.teammates[data.index].duties, 10);
        if (isNaN(duties)) {
          this.$refs.alertTeammate.showAlert("warning", "Warning", "Teammate duties must be number");
          return
        }

        // validate name must not be empty
        if (this.teammates[data.index].name == "") {
          this.$refs.alertTeammate.showAlert("warning", "Warning", "Teammate name must not be empty");
          return
        }

        // id == 0 --> create teammate
        if (this.teammates[data.index].id == 0) {
          // create teammate
          axios
              .post("api/teams/" + this.teamId + "/teammates",
                  {
                    "name": this.teammates[data.index].name,
                    "duty_readiness": this.teammates[data.index].duty_readiness,
                    "duties": duties,
                  },
                  { withCredentials: true })
              .then(
                  response => {
                    let teammateId = response.data.id;
                    this.teammates[data.index].id = teammateId;
                    this.showDuties = true;
                    this.$refs.alertTeammate.showAlert("success", "teammate '" + this.teammates[data.index].name + "' is created", "Success");
                  },
                  error => {
                    console.log('Error create teammate: ', error);
                    this.$refs.alertTeammate.showAlert("error", error, "Error");
                  },
              );
        } else {
          // update teammate
          axios
              .put("api/teammates/" + this.teammates[data.index].id,
                  {
                    "name": this.teammates[data.index].name,
                    "duty_readiness": this.teammates[data.index].duty_readiness,
                    "duties": duties,
                  },
                  { withCredentials: true })
              .then(
                  response => {
                    let status = response.data.status;
                    this.$refs.alertTeammate.showAlert("success", status, "Success");
                  },
                  error => {
                    console.log('Error update teammate: ', error);
                    this.$refs.alertTeammate.showAlert("error", error, "Error");
                  },
              );
        }
      }
    },
    addTeammate() {
      this.teammates.push(
          { id: 0, team_id: this.teamId, name: "", duty_readiness: true, duties: 0, isEdit: "true" }
      )
    },
    deleteTeammateRowHandler(data) {
      // clear empty teammates row (with id=0)
      if (this.teammates[data.index].id == 0) {

        let i = this.teammates.map(item => item.id).indexOf(0);
        this.teammates.splice(i, 1) // remove it from array
        return
      }

      this.dataToDelete = this.teammates[data.index];
      this.showModalTeammate = true;
    },
    cancelDeleteTeammate() {
      this.showModalTeammate = false;
      this.dataToDelete = {};
    },
    deleteTeammate() {
      // delete teammate
      if (this.dataToDelete != {}) {
        axios
            .delete("api/teammates/" + this.dataToDelete.id,
                { withCredentials: true })
            .then(
                response => {
                  let status = response.data.status;

                  let i = this.teammates.map(item => item.id).indexOf(this.dataToDelete.id);
                  if (i != -1) {
                    this.teammates.splice(i, 1)
                    this.$refs.alertTeammate.showAlert("success", status, "Success");

                    if (this.teammates.length == 0) {
                      this.showDuties = false;
                    }
                  } else {
                    // in case of team id is not found in this.teams
                    this.$refs.alertTeammate.showAlert("info", status + ", please, reload the page", "Info");
                  }

                  this.showModalTeammate = false;
                  this.dataToDelete = {};
                },
                error => {
                  console.log('Error delete teammate: ', error);
                  this.showModalTeammate = false;
                  this.dataToDelete = {};
                  this.$refs.alertTeammate.showAlert("error", error, "Error");
                },
            );
      }
    },
    changeDutyTeammate(teammateId, isDaily) {
      let dutyId;

      if (isDaily) {
        dutyId = this.dailyDutyId;
      } else {
        dutyId = this.weeklyDutyId;
      }
      console.log("isDaily", isDaily, "dutyId", dutyId, "teammateId", teammateId);

      if (dutyId != null) {
        // update duty
        axios
            .put("api/duties/" + dutyId,
                {
                  "teammate_id": teammateId,
                },
                { withCredentials: true })
            .then(
                response => {
                  if (isDaily) {
                    this.dailyDutyId = teammateId;
                  } else {
                    this.weeklyDutyId = teammateId;
                  }

                  this.$refs.alertTeammate.showAlert("success", response.data.status, "Success");
                },
                error => {
                  console.log('Error update duty: ', error);
                  this.$refs.alertTeammate.showAlert("error", error, "Error");
                },
            );
        } else {
        // create duty
        axios
            .post("api/teams/" + this.teamId + "/duties",
                {
                  "teammate_id": teammateId,
                  "is_daily": isDaily,
                },
                { withCredentials: true })
            .then(
                response => {
                  if (isDaily) {
                    this.dailyDutyId = response.data.id;
                  } else {
                    this.weeklyDutyId = response.data.id;
                  }

                  this.$refs.alertTeammate.showAlert("success", response.data.status, "Success");
                },
                error => {
                  console.log('Error create duty: ', error);
                  this.$refs.alertTeammate.showAlert("error", error, "Error");
                },
            );
      }
    },
    deleteDuty(isDaily) {
      // delete duty
      let dutyId;
      if (isDaily) {
        dutyId = this.dailyDutyId;
      } else {
        dutyId = this.weeklyDutyId;
      }

      if (dutyId == null) {
        this.$refs.alertTeammate.showAlert("info",  "duty is empty - can not be delete", "Info");
        return;
      }

        axios
            .delete("api/duties/" + dutyId,
                {withCredentials: true})
            .then(
                response => {
                  let status = response.data.status;
                  this.$refs.alertTeammate.showAlert("success", status, "Success");

                  if (isDaily) {
                    this.dailyDutyId = null;
                    this.selectedDailyTeammateId = null;
                  } else {
                    this.weeklyDutyId = null;
                    this.selectedWeeklyTeammateId = null;
                  }
                },
                error => {
                  console.log('Error delete duty: ', error);
                  this.$refs.alertTeammate.showAlert("error", error, "Error");
                },
            );
    }
  },
};
</script>