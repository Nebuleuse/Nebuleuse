angular.module('RDash')
	.controller('BuildCreateModal', ['$scope', '$http', '$uibModalInstance', 'list', 'commit', BuildCreateModal]);

function BuildCreateModal($scope, $http, $uibModalInstance, list, commit) {
	$scope.commits = [];
	$scope.diffs = [];
	$scope.list = list;
	$scope.fromCommit = commit;
	$scope.rawSize = 0;
	$http.post(APIURL + '/prepareGitBuild', {sessionid: $scope.Self.SessionId, commit: commit.Id})
		.success(function (data) {
			$scope.diffs = data.Diffs;
			$scope.rawSize = data.TotalSize;
		}).error(function (data, status) {
			$scope.parseError(data, status);
			$scope.addAlert("Could not fetch build preperation infos!", "danger");
		});

	var found = false;
	for (var i = 0; i < list.Commits.length; i++) {
		if (list.Commits[i].Id == commit.Id){
			found = true;
		}
		if (found){
			$scope.commits.push(list.Commits[i]);
		}
	};
}