import moment from 'moment';

moment.fn.toJSON = function () {
  return this.toISOString(true);
};

export default moment;
