function NodeCtrl($scope, $http) {
  $scope.nodes = [];
  $scope.working = false;

  var logError = function(data, status) {
    console.log('code '+status+': '+data);
    $scope.working = false;
  };

  var refresh = function() {
    return $http.get('/node/').
      success(function(data) { $scope.nodes = data.Nodes; }).
      error(logError);
  };

  $scope.addNode = function() {
    $scope.working = true;
    $http.post('/node/', {Name: $scope.nodeName, Role: $scope.nodeRole}).
      error(logError).
      success(function() {
        refresh().then(function() {
          $scope.working = false;
          $scope.nodeName = '';
          $scope.nodeRole = '';
        })
      });
  };

/*
  $scope.toggleDone = function(task) {
    data = {ID: task.ID, Title: task.Title, Done: !task.Done}
    $http.put('/task/'+task.ID, data).
      error(logError).
      success(function() { task.Done = !task.Done });
  };
*/

  refresh().then(function() { $scope.working = false; });
}
