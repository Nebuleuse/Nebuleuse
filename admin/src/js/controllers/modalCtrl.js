angular.module('RDash')
  .controller('ModalCtrl', ['$scope', '$modalInstance', ModalCtrl]);

function ModalCtrl($scope, $modalInstance) {
  $scope.ok = function () {
    $modalInstance.close();
  };

  $scope.cancel = function () {
    $modalInstance.dismiss('cancel');
  };
}