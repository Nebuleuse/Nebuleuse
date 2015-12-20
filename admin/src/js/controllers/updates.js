angular.module('RDash')
	.controller('UpdatesCtrl', ['$scope', '$http','$uibModal', UpdatesCtrl]);

function UpdatesCtrl($scope, $http, $uibModal) {
	$scope.setPageTitle("Update list");
	if(!$scope.checkAccess())
		return;
	$scope.list = [];
	$scope.fullList = [];

	$scope.refreshList = function () {
		$http.post(APIURL + '/getUpdateListWithGit', {sessionid: $scope.Self.SessionId, diffs: true})
		.success(function (data) {
			$scope.fullList = JSON.parse(JSON.stringify(data));
			var compare = "";
			if (data.Updates.length == 0){
				compare = data.CurrentCommit;
			} else {
				compare = data.Updates[0].Commit;
			}

			for (var i = data.Commits.length - 1; i >= 0; i--) {
				if (data.Commits[i].Id == compare){
					data.Commits = data.Commits.slice(0, i);
				}
			$scope.list = data;
			};
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

	$scope.refreshList();
}