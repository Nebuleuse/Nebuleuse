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
	$scope.setSelectedBranch  = function (obj) {
		$scope.selected = obj;
		$scope.selectedTpl = "templates/updates/branch.html";
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
	$scope.isRankAuth = function(accessRank, rank){
		return (accessRank & (1<<rank)) == (1<<rank);
	}
	$scope.toggleFiles = function(){
		$scope.showFiles = !$scope.showFiles;
	}
	$scope.refreshList = function () {
		$http.post(APIURL + '/getCompleteBranchUpdates', {sessionid: $scope.Self.SessionId, diffs: true})
		.success(function (data) {
			$scope.list = data;
			if (data.Builds.length == 0 && data.Commits.length > 0){
				$scope.addAlert("Looks like you do not have any builds yet. Select a commit to create a build.", "info");
			}
			if(data.Builds !== null){
				for	(var i=0; i < data.Builds.length; i++){
					if(data.Builds[i].FileChanged === ""){
						continue;
					}
					$scope.list.Builds[i].FileChanged = JSON.parse(data.Builds[i].FileChanged);
				}
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

    $scope.createEmptyBranch = function(){
        var modalInstance = $uibModal.open({
	      animation: true,
	      templateUrl: 'templates/updates/createBranchModal.html',
	      controller: 'BranchCreateModal',
	      scope: $scope,
	      size: 'lg',
		  resolve: {
			  build: function(){return 0;}
		  }
	    });
    }
    $scope.createBranchFromBuild = function(build){
        var modalInstance = $uibModal.open({
	      animation: true,
	      templateUrl: 'templates/updates/createBranchModal.html',
	      controller: 'BranchCreateModal',
	      scope: $scope,
	      size: 'lg',
		  resolve: {
			  build: function(){return build.Id;}
		  }
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
	$scope.downloadUpdate = function(branch, buildId){
		var branchObj;
		for(var j = 0; j < $scope.list.Branches.length; j++){
			if($scope.list.Branches[j].Name == branch){
				branchObj = $scope.list.Branches[j];
				break;
			}
		}
		var updates = branchObj.Updates;
		var i = 0;
		for (i=0; i < updates.length; i++){
			if(updates[i].BuildId == buildId){
				break;
			}
		}
		var from = 0;
		if(i+1 < updates.length)
			from = updates[i+1].BuildId;

		var url = APIURL + "/" + $scope.Nebuleuse.UpdatesLocation + from + "to" + buildId + ".tar.xz";
		window.open(url);
	}
	$scope.refreshList();
}