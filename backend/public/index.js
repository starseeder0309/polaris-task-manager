(function ($) {
  'use strict';
  $(function () {
    var taskList = $('.task-list');
    var taskListInput = $('.task-list-input');

    $('.task-list-add-button').on('click', function (event) {
      event.preventDefault();

      var item = $(this).prevAll('.task-list-input').val();

      if (item) {
        $.post('/tasks', JSON.stringify({ title: item }), addItem);
        taskListInput.val('');
      }
    });

    var addItem = function (item) {
      if (item.isCompleted) {
        taskList.append(
          "<li class='completed'" +
            " id='" +
            item.id +
            "'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' checked='checked' />" +
            item.title +
            "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>"
        );
      } else {
        taskList.append(
          '<li ' +
            " id='" +
            item.id +
            "'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" +
            item.title +
            "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>"
        );
      }
    };

    $.get('/tasks', function (items) {
      items.forEach((e) => {
        addItem(e);
      });
    });

    taskList.on('change', '.checkbox', function () {
      var id = parseInt($(this).closest('li').attr('id'));
      var title = $(this).closest('li').text();
      var $self = $(this);
      var complete = true;
      if ($(this).attr('checked')) {
        complete = false;
      }
      $.ajax({
        url: '/tasks/' + id,
        type: 'PUT',
        data: JSON.stringify({ id: id, title: title, isCompleted: complete }),
        success: function (data) {
          if (complete) {
            $self.attr('checked', 'checked');
          } else {
            $self.removeAttr('checked');
          }

          $self.closest('li').toggleClass('completed');
        },
      });
    });

    taskList.on('click', '.remove', function () {
      var id = $(this).closest('li').attr('id');
      var $self = $(this);
      $.ajax({
        url: '/tasks/' + id,
        type: 'DELETE',
        success: function (data) {
          if (data.success) {
            $self.parent().remove();
          }
        },
      });
    });
  });
})(jQuery);
