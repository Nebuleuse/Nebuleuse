'use strict';

/**
 * Route configuration for the RDash module.
 */
angular.module('RDash').config(['$stateProvider', '$urlRouterProvider',
    function($stateProvider, $urlRouterProvider) {

        // For unmatched routes
        $urlRouterProvider.otherwise('/');

        // Application routes
        $stateProvider
            .state('dashboard', {
                url: '/',
                controller: 'DashboardCtrl',
                templateUrl: 'templates/dashboard.html'
            })
            .state('login', {
                url: '/login',
                controller: 'LoginCtrl',
                templateUrl: 'templates/login.html'
            })
            .state('noauth', {
                url: '/noauth',
                templateUrl: 'templates/noauth.html'
            })
            .state('log', {
                url: '/log',
                controller: 'LogCtrl',
                templateUrl: 'templates/log.html'
            })
            .state('users', {
                url: '/users',
                controller: 'UserListCtrl',
                templateUrl: 'templates/usersList.html'
            })
            .state('user', {
                url: '/user/:userId',
                controller: 'UserCtrl',
                templateUrl: 'templates/user.html'
            })
            .state('achievements', {
                url: '/achievements',
                controller: 'AchievementsCtrl',
                templateUrl: 'templates/achievementsList.html'
            })
            .state('achievementEdit', {
                url: '/achievementEdit/:achievementId',
                controller: 'AchievementCtrl',
                templateUrl: 'templates/achievementEdit.html'
            })
            .state('achievementAdd', {
                url: '/achievementAdd',
                controller: 'AchievementCtrl',
                templateUrl: 'templates/achievementEdit.html'
            })
            .state('stats', {
                url: '/stats',
                controller: 'StatsCtrl',
                templateUrl: 'templates/statsList.html'
            })
            .state('statTableEdit', {
                url: '/statTableEdit/:statName',
                controller: 'StatTableEditCtrl',
                templateUrl: 'templates/statsEdit.html'
            })
            .state('statTableAdd', {
                url: '/statTableAdd',
                controller: 'StatTableEditCtrl',
                templateUrl: 'templates/statsEdit.html'
            })
            .state('matchmaking', {
                url: '/matchmaking',
                templateUrl: 'templates/wip.html'
            })
            .state('servers', {
                url: '/servers',
                templateUrl: 'templates/wip.html'
            })
            .state('items', {
                url: '/items',
                templateUrl: 'templates/wip.html'
            })
            .state('config', {
                url: '/config',
                templateUrl: 'templates/wip.html'
            });
    }
]);