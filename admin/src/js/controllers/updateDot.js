angular.module('RDash')
	.controller('updateDotCtrl', ['$scope', updateDotCtrl]);

function updateDotCtrl($scope) {
    var updates = $scope.branch.Updates;
    $scope.update = {};
    $scope.exist = false;
    $scope.visible = true;
	for (var i=0; i < updates.length; i++){
		if(updates[i].BuildId == $scope.build.Id){
			$scope.update = updates[i];
            $scope.exist = true;
		}
	}
    if(!$scope.exist){
        if($scope.branch.Updates[0].BuildId >= $scope.build.Id){
            $scope.visible = false;
        }
    }
}