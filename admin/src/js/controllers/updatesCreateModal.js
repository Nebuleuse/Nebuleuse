angular.module('RDash')
	.controller('UpdateCreateModal', ['$scope', '$http', '$uibModalInstance', 'list', 'commit', UpdateCreateModal]);

function UpdateCreateModal($scope, $http, $uibModalInstance, list, commit) {
	$scope.commits = [];
	$scope.diffs = [];
	$scope.list = list;
	$scope.fromCommit = commit;
	$http.post(APIURL + '/prepareGitPatch', {sessionid: $scope.Self.SessionId, commit: commit})
		.success(function (data) {
			console.log(data)
			$scope.diffs = data.Diffs;
			$scope.rawSize = data.TotalSize;
		});
	var compare = "";
	if (list.Updates.length == 0){
		compare = list.CurrentCommit;
	} else {
		compare = list.Updates[0].Commit;
	}
	var found = false;
	for (var i = 0; i < list.Commits.length; i++) {
		if (list.Commits[i].Id == commit){
			found = true;
		}
		if (list.Commits[i].Id == compare){
			break;
		}
		if (found){
			$scope.commits.push(list.Commits[i]);
		}
	};
}