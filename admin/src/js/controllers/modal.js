angular.module('RDash')
  .controller('ModalCtrl', ['$scope', '$uibModalInstance', ModalCtrl]);

function ModalCtrl($scope, $uibModalInstance) {
  $scope.ok = function () {
    $uibModalInstance.close();
  };

  $scope.cancel = function () {
    $uibModalInstance.dismiss('cancel');
  };
}