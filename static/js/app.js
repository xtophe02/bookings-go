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
    showConfirmButton = true
  }) => {
    Swal.fire({
      title,
      html,
      icon,
      confirmButtonText,
      showConfirmButton
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

        const rangePicker = new DateRangePicker(elem, {
          minDate: new Date(),
          format: 'yyyy-mm-dd',
          showOnFocus: true,
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

const reservationDates =(roomID,csrf) =>{
  let attention = Prompt();
  const submit = document.querySelector('#submit');

  submit.addEventListener('click', () => {
    let html = `
  <form class="row needs-validation mx-auto" novalidate id="reservations-dates" style="max-width: 400px;" autocomplete="off">
    <div class="col">
      <input id='start' type="text" required class="form-control" name="start_date" placeholder="Arrival" />
    </div>
    <div class="col">
      <input id='end' type="text" required class="form-control" name="end_date" placeholder="Departure" />
    </div>
  </form>
  `;
    attention.form({
      title: 'Chose your dates',
      html,
      cb: async (res) => {
        const form = document.querySelector('#reservations-dates');
        const formData = new FormData(form);
        formData.append('csrf_token', csrf);
        formData.append("room_id",roomID)
      
        const response = await fetch('/availability-json', {
          method: 'POST',
          body: formData,
       
        });

        const json = await response.json();
        if (json.ok){
          attention.modal({
            title: "Room is available!",
            html: `<a class="btn btn-primary my-2" href="/book-room?id=${json.room_id}&s=${json.start_date}&e=${json.end_date}">Book Now</a>`,
            icon: 'success',
            showConfirmButton:false
          })
        }else{
          attention.modal({
            title: "No Availability",
            html: '<em>Please to enter new dates</em>',
            icon: 'error',
            confirmButtonText :'ok',
          })
        }
      },
    });
  });
}