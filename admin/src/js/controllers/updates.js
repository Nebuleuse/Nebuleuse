angular.module('RDash')
	.controller('UpdatesCtrl', ['$scope', '$http','$uibModal', UpdatesCtrl]);

function UpdatesCtrl($scope, $http, $uibModal) {
	$scope.setPageTitle("Update list");
	if(!$scope.checkAccess())
		return;
	$scope.list = {};
	$scope.selected = {};
	$scope.selectedTpl = "";
	
	$scope.setSelectedCommit  = function (obj) {
		$scope.selected = obj;
		$scope.selectedTpl = "templates/updates/commits.html";
	}
	$scope.setSelectedBuild  = function (obj) {
		$scope.selected = obj;
		$scope.selectedTpl = "templates/updates/builds.html";
	}
	$scope.setSelectedUpdate  = function (obj) {
		$scope.selected = obj;
		$scope.selectedTpl = "templates/updates/updates.html";
	}
	
	$scope.refreshList = function () {
		$http.post(APIURL + '/getCompleteBranchUpdates', {sessionid: $scope.Self.SessionId, diffs: true})
		.success(function (data) {
			console.log(data)
			$scope.list = data;
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch updates infos!", "danger");
		});
	}

	$scope.updateCacheList = function () {
		$http.post(APIURL + '/updateGitCommitCacheList', {sessionid: $scope.Self.SessionId}).success(function (data) {
			$scope.refreshList();
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not update cache list!", "danger");
		})
	}

	$scope.createPatch = function (commit) {
		var modalInstance = $uibModal.open({
	      animation: true,
	      templateUrl: 'templates/updates/createModal.html',
	      controller: 'UpdateCreateModal',
	      scope: $scope,
	      size: 'lg',
	      resolve: {
	        list: function () {
	          return $scope.fullList;
	        },
	        commit: function(){
	        	return commit;
	        }
	      }
	    });
	}
	$scope.getUpdateForBuild = function (branch, id) {
		var updates = branch.Updates;
		for (var i=0; i < updates.length; i++){
			if(updates[i].BuildId == id){
				return updates[i];
			}
		}
		return {};
	}
	$scope.refreshList();
}