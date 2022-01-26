import moment from './moment';

test('moment.toJSON()', () => {
  const date = moment('2021-12-12T17:34:00-05:00');
  expect(date.toJSON()).toEqual('2021-12-12T17:34:00.000-05:00');
  expect(JSON.stringify(date)).toEqual('"2021-12-12T17:34:00.000-05:00"');
  expect(JSON.stringify({ date: date })).toEqual('{"date":"2021-12-12T17:34:00.000-05:00"}');
});
