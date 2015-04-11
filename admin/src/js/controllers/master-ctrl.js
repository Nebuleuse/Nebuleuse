/**
 * Master Controller
 */

angular.module('RDash')
    .controller('MasterCtrl', ['$scope', '$cookieStore','$http', '$location', '$rootScope', MasterCtrl]);

function MasterCtrl($scope, $cookieStore, $http, $location, $rootScope) {
    $scope.alerts = [];

    $scope.addAlert = function(message, level) {
        $scope.alerts.push({
            msg: message,
            type: level
        });
    };

    $scope.closeAlert = function(index) {
        $scope.alerts.splice(index, 1);
    };

   $scope.getUserInfos = function () {
        $http.post(APIURL + '/getUserInfos', {sessionid: $scope.Self.SessionId, infomask:UserMaskBase})
        .success(function (data) {
            $scope.Self = data;
            if(data.Rank < 2){
                $location.path('/noauth');
                $scope.setConnected(false);
            } else {
                $scope.setConnected(true);
            }
        }).error(function (data, status) {
            $scope.setConnected(false);
        });
    };

    $scope.setPageTitle = function (title) {
        $rootScope.PageTitle = title + " - Nebuleuse";
        $scope.PageTitle = title;
    }
    $scope.setConnected = function(connected) {
        if(!connected) {
            $scope.Self = {};
            $cookieStore.remove('sessionId');
            $location.path('/login');
        }
        $scope.isConnected = connected;
    };
    $scope.checkAccess = function() {
        if(!$scope.isConnected || $scope.Self.Rank < 2){
            $location.path('/login');
            return false;
        }
        return true;
    };

    $scope.Menus = [    {name: "Home", icon: "fa-home", link:"/"},
                        {name: "Live", icon: "fa-cloud", link:"live"},
                        {name: "Achievements", icon: "fa-trophy", link:"achievements"},
                        {name: "Stats", icon: "fa-pie-chart", link:"stats"},
                        {name: "Matchmaking", icon: "fa-globe", link:"matchmaking"},
                        {name: "Users", icon: "fa-users", link:"users"},
                        {name: "Servers", icon: "fa-server", link:"servers"},
                        {name: "Items", icon: "fa-sitemap", link:"items"}];
    $scope.setPageTitle("Dashboard");
    $scope.isConnected = false;
    $scope.Self = {};
    
    if(angular.isDefined($cookieStore.get('sessionId'))){
        $scope.Self.SessionId = $cookieStore.get('sessionId');
        $scope.getUserInfos();
    } else if ($location.path() != '/login') {
        $scope.setConnected(false);
    }
    

    $scope.logout = function () {
        $http.post(APIURL + '/disconnect', {sessionid: $scope.Self.SessionId});
        $scope.setConnected(false);
    };
    
    var mobileView = 992;

    $scope.getWidth = function() {
        return window.innerWidth;
    };

    $scope.$watch($scope.getWidth, function(newValue, oldValue) {
        if (newValue >= mobileView) {
            if (angular.isDefined($cookieStore.get('toggle'))) {
                $scope.toggle = ! $cookieStore.get('toggle') ? false : true;
            } else {
                $scope.toggle = true;
            }
        } else {
            $scope.toggle = false;
        }

    });

    $scope.toggleSidebar = function() {
        $scope.toggle = !$scope.toggle;
        $cookieStore.put('toggle', $scope.toggle);
    };

    window.onresize = function() {
        $scope.$apply();
    };
}