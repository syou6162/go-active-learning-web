import moment from 'moment'

export default function NewExample(e) {
  var yesterday = moment().add(-1, "days");
  var createdAt = moment(e.CreatedAt);
  var updatedAt = moment(e.UpdatedAt);
  e.CreatedAt = createdAt;
  e.UpdatedAt = updatedAt;
  e.IsNew = createdAt.isAfter(yesterday);
  return e;
}
