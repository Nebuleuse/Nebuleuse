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

    $scope.parseError = function(data, status) {
        if(data.Code == 0)
            return;
        if(data.Code == 2 || status == 401)
            $scope.setConnected(false);
        console.error(data, status);
    };
    $scope.parseMessage = function(data) {
        if (data.Code != null && data.Code == 0)
            return; // Longpoll timedout

        $scope.$broadcast(data.Channel, data.Message);
    };
    $scope.getUserInfos = function () {
        $http.post(APIURL + '/getUserInfos', {sessionid: $scope.Self.SessionId, infomask:UserMaskBase})
        .success(function (data) {
            var sessionid = $scope.Self.SessionId;
            $scope.Self = data;
            $scope.Self.SessionId = sessionid;
            if(data.Rank < 2){
                $scope.setConnected(false);
                $location.path('/noauth');
            } else {
                $scope.setConnected(true);
                if($location.path() === "/login") {
                    $location.path('/');
                }
            }
        }).error(function (data, status) {
            $scope.parseError(data, status);
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
    $scope.subscribeTo = function (channel) {
        $http.post(APIURL + '/subscribeTo', {sessionid: $scope.Self.SessionId, channel:channel})
        .error(function (data, status) {
            $scope.parseError(data, status);
            $scope.addAlert("Could not subscribe to " + channel, "danger");
        });
    }
    $scope.unSubscribeTo = function (channel) {
        $http.post(APIURL + '/unSubscribeTo', {sessionid: $scope.Self.SessionId, channel:channel})
        .error(function (data, status) {
            $scope.parseError(data, status);
            $scope.addAlert("Could not unsubscribe to " + channel, "danger");
        });
    }
    $scope.getMessages = function () {
        $scope.lostConnection = false;
        $http.post(APIURL + '/getMessages', {sessionid: $scope.Self.SessionId})
        .success(function (data) {
            $scope.parseMessage(data);
            $scope.getMessages();
        })
        .error(function (data, status) {
            if (status === -1) // connection error results in -1
                return;
            $scope.parseError(data, status);
            if($scope.isConnected)
                $scope.addAlert("Could not get new messages", "danger");
            $scope.lostConnection = true;
        });
    }

    $scope.goto = function (path) {
        $location.path(path);
    }

    $scope.Menus = [    {name: "Home", icon: "fa-home", link:"/"},
                        {name: "Log", icon: "fa-cloud", link:"log"},
                        {name: "Users", icon: "fa-users", link:"users"},
                        {name: "Achievements", icon: "fa-trophy", link:"achievements"},
                        {name: "Stats", icon: "fa-pie-chart", link:"stats"},
                        {name: "Updates", icon: "fa-folder-open", link:"updates"},
                        {name: "Matchmaking", icon: "fa-globe", link:"matchmaking"},
                        {name: "Servers", icon: "fa-server", link:"servers"},
                        {name: "Items", icon: "fa-sitemap", link:"items"}];
    $scope.setPageTitle("Dashboard");
    $scope.isConnected = false;
    $scope.lostConnection = false;
    $scope.Self = {};
    $scope.Nebuleuse = {};
    
    $http.get(APIURL + '/status')
        .success(function (data) {
            $scope.Nebuleuse = data;
            if(angular.isDefined($cookieStore.get('sessionId'))){
                $scope.Self.SessionId = $cookieStore.get('sessionId');
                $scope.getUserInfos();
                $scope.getMessages();
            } else if ($location.path() != '/login') {
                $scope.setConnected(false);
            }
        }).error(function (data, status) {
            if (status === -1)
                return;
            $scope.parseError(data, status);
            $scope.addAlert("Could not get server status", "danger");
        })
    

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