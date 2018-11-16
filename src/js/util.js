var moment = require('moment-timezone');

export default function NewExample(e, opts = {}) {
  var isNewDayThreshold = moment().add(-1 * (opts["IsNewDayThreshold"] || 1), "days");
  var createdAt = moment(e.CreatedAt);
  var updatedAt = moment(e.UpdatedAt);
  e.CreatedAt = createdAt;
  e.UpdatedAt = updatedAt;
  e.IsNew = createdAt.isAfter(isNewDayThreshold);
  return e;
}
