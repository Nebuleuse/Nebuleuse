angular.module('RDash')
	.controller('UpdatesCtrl', ['$scope', '$http','$uibModal', UpdatesCtrl]);

function UpdatesCtrl($scope, $http, $uibModal) {
	$scope.setPageTitle("Update list");
	if(!$scope.checkAccess())
		return;
	$scope.list = {};
	$scope.selected = {};
	$scope.selectedBranch={};
	$scope.selectedTpl = "";
	$scope.showFiles = false;
	
	$scope.setSelectedCommit  = function (obj) {
		$scope.selected = obj;
		$scope.selectedTpl = "templates/updates/commits.html";
	}
	$scope.setSelectedBuild  = function (obj) {
		$scope.selected = obj;
		$scope.selectedTpl = "templates/updates/builds.html";
	}
	$scope.setSelectedUpdate  = function (obj, branch) {
		if (obj.create !== null && obj.create === true){
			$scope.createPatch(obj.build, obj.branch);
		} else {
			$scope.selected = obj;
			$scope.selectedBranch = branch;
			$scope.selectedTpl = "templates/updates/updates.html";
		}
	}
	$scope.toggleFiles = function(){
		$scope.showFiles = !$scope.showFiles;
	}
	$scope.refreshList = function () {
		$http.post(APIURL + '/getCompleteBranchUpdates', {sessionid: $scope.Self.SessionId, diffs: true})
		.success(function (data) {
			$scope.list = data;
			console.log($scope.list);
			if (data.Builds.length == 0 && data.Commits.length > 0){
				$scope.addAlert("Looks like you do not have any builds yet. Select a commit to create a build.", "info");
			}
			for	(var i=0; i < data.Builds.length; i++){
				if(data.Builds[i].FileChanged === ""){
					continue;
				}
				$scope.list.Builds[i].FileChanged = JSON.parse(data.Builds[i].FileChanged);
			}
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
	$scope.createBuild = function (commit) {
		var modalInstance = $uibModal.open({
	      animation: true,
	      templateUrl: 'templates/updates/createBuildModal.html',
	      controller: 'BuildCreateModal',
	      scope: $scope,
	      size: 'lg',
	      resolve: {
	        list: function () {
	          return $scope.list;
	        },
	        commit: function(){
	        	return commit;
	        }
	      }
	    });
	}
	$scope.createPatch = function (build, branch) {
		var modalInstance = $uibModal.open({
	      animation: true,
	      templateUrl: 'templates/updates/createUpdateModal.html',
	      controller: 'UpdateCreateModal',
	      scope: $scope,
	      size: 'lg',
	      resolve: {
			  list: function () {
				 return $scope.list;
			  },
	        build: function () {
	          return build;
	        },
	        branch: function(){
	        	return branch;
	        }
	      }
	    });
	}
	$scope.createManualBuild = function(){
		var modalInstance = $uibModal.open({
	      animation: true,
	      templateUrl: 'templates/updates/createManualBuildModal.html',
	      controller: 'ManualBulidCreateModal',
	      scope: $scope,
	      size: 'lg'
	    });
	}
	$scope.setActiveUpdate = function (update, branch) {
		$http.post(APIURL + '/setActiveUpdate', {sessionid: $scope.Self.SessionId, build: update.BuildId, branch: branch.Name})
		.success(function () {
			$scope.refreshList();
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not set active update!", "danger");
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
	$scope.downloadUpdate = function(from, to){
		var url = APIURL + "/" + $scope.Nebuleuse.UpdatesLocation + from + "to" + to + ".tar.xz";
		console.log(url) 
	}
	$scope.refreshList();
}