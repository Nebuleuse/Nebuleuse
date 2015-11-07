angular.module('RDash')
.directive("updateList", function(){
  return {
    restrict: "A",
    scope: {
      list: '=updateList'
    },
    link: function(scope, element, attrs){
      var ctx = element[0].getContext('2d');

      function updateCanvas(data) {
        draw(0, 10, element[0].width, 10)
      }

      scope.$watch(scope.list, function(value) {
        console.dir(value)
        updateCanvas(value)
      });
        console.dir(scope.list)
        console.dir(scope)

      element.bind('mousedown', function(event){
      });
      element.bind('mousemove', function(event){
      });
      element.bind('mouseup', function(event){
      });

      // canvas reset
      function reset(){
       element[0].width = element[0].width; 
      }

      function draw(lX, lY, cX, cY){
        // line from
        ctx.moveTo(lX,lY);
        // to
        ctx.lineTo(cX,cY);
        // color
        ctx.strokeStyle = "#4bf";
        // draw it
        ctx.stroke();
      }
    }
  };
});