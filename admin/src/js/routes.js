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
            .state('achievements', {
                url: '/achievements',
                controller: 'AchievementsCtrl',
                templateUrl: 'templates/achievements.html'
            })
            .state('stats', {
                url: '/stats',
                templateUrl: 'templates/stats.html'
            })
            .state('matchmking', {
                url: '/matchmking',
                templateUrl: 'templates/matchmking.html'
            })
            .state('users', {
                url: '/users',
                templateUrl: 'templates/users.html'
            })
            .state('servers', {
                url: '/servers',
                templateUrl: 'templates/servers.html'
            })
            .state('items', {
                url: '/items',
                templateUrl: 'templates/items.html'
            })
            .state('config', {
                url: '/config',
                templateUrl: 'templates/config.html'
            });
    }
]);