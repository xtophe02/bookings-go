const Prompt = () => {
  const toast = ({ title = '', icon = 'success', position = 'top-end' }) => {
    const Toast = Swal.mixin({
      toast: true,
      position,
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      didOpen: (toast) => {
        toast.addEventListener('mouseenter', Swal.stopTimer);
        toast.addEventListener('mouseleave', Swal.resumeTimer);
      },
    });

    Toast.fire({
      icon,
      title,
    });
  };
  const modal = ({
    title = '',
    html,
    icon = 'success',
    confirmButtonText = 'Cool',
  }) => {
    Swal.fire({
      title,
      html,
      icon,
      confirmButtonText,
    });
  };
  const form = async ({ title, html, cb }) => {
    const { value: result } = await Swal.fire({
      title,
      html,
      backdrop: false,
      focusConfirm: false,
      showCancelButton: true,
      willOpen: () => {
        const elem = document.getElementById('reservations-dates');

        new DateRangePicker(elem, {
          format: 'yyyy-mm-dd',
          showOnFocus: false,
        });
      },
      preConfirm: () => {
        return [
          document.getElementById('start').value,
          document.getElementById('end').value,
        ];
      },
    });

    if (result) {
      if (result.dismiss !== Swal.DismissReason.cancel && result.value !== '') {
        cb(result);
      }
      // Swal.fire(JSON.stringify(formValues));
    }
  };
  const notify = (text, type) => {
    notie.alert({ type, text });
  };
  return { toast, modal, form, notify };
};
