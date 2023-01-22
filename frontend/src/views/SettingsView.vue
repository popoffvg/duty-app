<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-8 offset-sm-2">
        <div>

          <div class="table-header">
            <h2>Team  "{{ team.title }}" settings</h2>
          </div>

            <div class="settings-container">
              <label for="space-channel">Space channel name (case sensitive, space app must be authorize in channel)</label>
              <input type="text"
                     id="space-channel"
                     class="form-control"
                     v-model="team.space_channel">
            </div>

          <div class="buttons-panel">
            <b-button class="save-button" @click="saveChannelName()">
              <span>Save</span>
            </b-button>

            <b-button v-if="team.space_channel" class="test-alert-button" @click="sendTestNotification()">
              <span>Send test message</span>
            </b-button>
          </div>

        </div>
      </div>
    </div>
  </div>

  <vue-basic-alert :duration="300" :closeIn="3000" ref="alertSettings" />
</template>

<style>
.buttons-panel {
  padding-top: 10px;
}

.buttons-panel button {
  margin-right: 10px;
}

.save-button {
  background-color: green !important;
}
</style>

<script>
import axios from "axios";
import VueBasicAlert from 'vue-basic-alert'

export default {
  components: {
    VueBasicAlert,
  },
  data() {
    return {
      team: {
        id: this.$route.params.teamId,
        title: "my-team",
        space_channel: "duty",
      },
    };
  },
  created() {
    axios
        .get("api/teams/" + this.team.id,{ withCredentials: true })
        .then(
            response => {
              if (response.data !== null) {
                this.team = response.data;
              } else {
                this.$refs.alertSettings.showAlert("error", "Can't get team data", "Error");
              }
            },
            error => {
              console.log('Error load team data: ', error);

              if (error.response.status === 401) {
                this.$router.push('/');
                return;
              }

              this.$refs.alertSettings.showAlert("error", error, "Error");
            },
        );
  },
  methods: {
    saveChannelName() {
      if (this.team.space_channel == "") {
        this.$refs.alertSettings.showAlert("warning", "Warning", "Channel name must not be empty");
        return
      }

      if (this.team.title == "") {
        this.$refs.alertSettings.showAlert("error", "Error", "Internal error");
        return
      }

      axios
          .put("api/teams/" + this.team.id,
              {
                "title": this.team.title,
                "space_channel": this.team.space_channel
              },
              { withCredentials: true })
          .then(
              response => {
                let status = response.data.status;
                this.$refs.alertSettings.showAlert("success", status, "Success");
              },
              error => {
                console.log('Error update team settings: ', error);
                this.$refs.alertSettings.showAlert("error", error, "Error");
              },
          );
    },
    sendTestNotification() {
      axios
          .post("api/teams/" + this.team.id + "/notifications/test",
              {},
              { withCredentials: true })
          .then(
              () => {
                this.$refs.alertSettings.showAlert("success", "notification was send", "Success");
              },
              error => {
                console.log('Error send test notification: ', error);
                this.$refs.alertSettings.showAlert("error", error, "Error");
              },
          );
    }
  },
};
</script>