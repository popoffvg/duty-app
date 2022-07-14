export default {
    formatDateDaily: function (dateString) {
        let date = new Date(dateString);
        let month = (1 + date.getMonth()).toString();
        month = month.length > 1 ? month : '0' + month;

        let day = date.getDate().toString();
        day = day.length > 1 ? day : '0' + day;

        return day + '/' + month;
    },
    formatDateWeekly: function (dateString) {
        let monday = new Date(dateString);

        // diff from monday to sunday: (seconds * minutes * hours * milliseconds = 1 day) * 6 day
        let diffToSunday = 518400000;
        let sunday = new Date(monday.getTime() + diffToSunday);

        let result = "";

        let mondayMonth = (1 + monday.getMonth()).toString();
        mondayMonth = mondayMonth.length > 1 ? mondayMonth : '0' + mondayMonth;

        let sundayMonth = (1 + sunday.getMonth()).toString();
        sundayMonth = sundayMonth.length > 1 ? sundayMonth : '0' + sundayMonth;

        let mondayDay = monday.getDate().toString();
        mondayDay = mondayDay.length > 1 ? mondayDay : '0' + mondayDay;

        let sundayDay = sunday.getDate().toString();
        sundayDay = sundayDay.length > 1 ? sundayDay : '0' + sundayDay;

        if (mondayMonth == sundayMonth) {
            result = mondayDay + "−" + sundayDay + "/" + mondayMonth;
        } else {
            result = mondayDay + "/" + mondayMonth + " − " + sundayDay + "/" + sundayMonth;
        }

        return result;
    }
}
