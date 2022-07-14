<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-8 offset-sm-2">
        <div>

          <div class="table-header">
            <h2>My Teams</h2>
          </div>

          <div class="add-button">
            <b-button @click="addTeam()">
              <span>Add team</span>
            </b-button>
          </div>

          <div class="teams-table">
            <b-table :items="teams" :fields="fields">
              <template #cell(index)="data">
                <span>{{ data.index + 1}}</span>
              </template>
              <template #cell(title)="data">
                <b-form-input v-if="teams[data.index].isEdit" type="text" v-model="teams[data.index].title"></b-form-input>
                <a v-else :href="'#/teams/' + teams[data.index].id + '/teammates'" >{{ data.value }}</a>
              </template>
              <template #cell(options)="data">
                <div class="table-buttons">
                  <b-button @click="editRowHandler(data)">
                    <span v-if="!teams[data.index].isEdit">Edit</span>
                    <span v-else>Save</span>
                  </b-button>
                  &nbsp;
                  <b-button @click="deleteTeamRowHandler(data)" class="delete-button">
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

  <div v-if="showModalTeam" class="modal">
    <!-- Modal content -->
    <div class="modal-content">
      <span class="close" @click="cancelDelete()">&times;</span>
      <p>Are you sure to delete team '{{ dataToDelete.title }}'?</p>

      <div class="modal-buttons">
        <b-button @click="cancelDelete()">
          <span>cancel</span>
        </b-button>
        <b-button class="delete-button" @click="deleteTeam()">
          <span>delete</span>
        </b-button>
      </div>
    </div>
  </div>


  <vue-basic-alert :duration="300" :closeIn="3000" ref="alertTeam" />
</template>

<script>
import axios from 'axios'
import VueBasicAlert from 'vue-basic-alert'


export default {
  components: {
    VueBasicAlert,
  },
  data() {
    return {
      showModalTeam: false,
      dataToDelete: {},
      fields: [
        { key: "index", label: "#"},
        { key: "title", label: "Title"},
        { key: "options", label: "Options"},
      ],
      teams: [
        // { id: 18, user_id: 2, title: 'Team New-1-2-3' },
        // { id: 19, user_id: 2, title: 'Team ABC' },
      ],
    };
  },
  created() {
    // load teams
     axios
         .get("api/teams",{ withCredentials: true })
         .then(
             response => {
               if (response.data.data !== null) {
                 this.teams = response.data.data;
                 this.teams = this.teams.map(item => ({...item, isEdit: false}));
               } else {
                 this.$refs.alertTeam.showAlert("info", "Click on [Add team] button ", "Info");
               }
             },
             error => {
               console.log('Error load teams: ', error);

               if (error.response.status === 401) {
                 this.$router.push('/');
                 return;
               }

               this.$refs.alertTeam.showAlert("error", error, "Error");
             },
         );
  },
  methods: {
    editRowHandler(data) {
      this.teams[data.index].isEdit = !this.teams[data.index].isEdit;

      if (!this.teams[data.index].isEdit) {
        if (this.teams[data.index].id == 0) {
          if (this.teams[data.index].title == "") {
            this.$refs.alertTeam.showAlert("warning", "Warning", "Team title must not be empty");
            return
          }

          // create team
          axios
              .post("api/teams/",
                  { "title": this.teams[data.index].title },
                  { withCredentials: true })
              .then(
                  response => {
                    let teamId = response.data.id;
                    this.teams[data.index].id = teamId;
                    this.$refs.alertTeam.showAlert("success", "team '" + this.teams[data.index].title + "' is created", "Success");
                  },
                  error => {
                    console.log('Error create team: ', error);
                    this.$refs.alertTeam.showAlert("error", error, "Error");
                  },
              );
        } else {
          // update team
          axios
              .put("api/teams/" + this.teams[data.index].id,
                  { "title": this.teams[data.index].title },
                  { withCredentials: true })
              .then(
                  response => {
                    let status = response.data.status;
                    this.$refs.alertTeam.showAlert("success", status, "Success");
                  },
                  error => {
                    console.log('Error update team: ', error);
                    this.$refs.alertTeam.showAlert("error", error, "Error");
                  },
              );
        }
      }
    },
    addTeam() {
      this.teams.push(
          { id: 0, user_id: this.teamId, title: "", isEdit: "true" }
      )
    },
    deleteTeamRowHandler(data) {
      // clear empty teams row (with id=0)
      if (this.teams[data.index].id == 0) {

        let i = this.teams.map(item => item.id).indexOf(0);
        this.teams.splice(i, 1) // remove it from array
        return
      }

      this.dataToDelete = this.teams[data.index];
      this.showModalTeam = true;
    },
    cancelDelete() {
      this.showModalTeam = false;
      this.dataToDelete = {};
    },
    deleteTeam() {
      // delete team
      if (this.dataToDelete != {}) {
        axios
            .delete("api/teams/" + this.dataToDelete.id,
                { withCredentials: true })
            .then(
                response => {
                  let status = response.data.status;

                  let i = this.teams.map(item => item.id).indexOf(this.dataToDelete.id);
                  if (i != -1) {
                    this.teams.splice(i, 1)
                    this.$refs.alertTeam.showAlert("success", status, "Success");
                  } else {
                    // in case of team id is not found in this.teams
                    this.$refs.alertTeam.showAlert("info", status + ", please, reload the page", "Info");
                  }

                  this.showModalTeam = false;
                  this.dataToDelete = {};
                },
                error => {
                  console.log('Error update team: ', error);
                  this.showModalTeam = false;
                  this.dataToDelete = {};
                  this.$refs.alertTeam.showAlert("error", error, "Error");
                },
            );
      }
    },
  },
};
</script>