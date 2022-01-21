import ServiceResponse from './ServiceResponse';

test('response.getKeyComponents()', () => {
  const response = new ServiceResponse(
    'ffbc9005-c62a-4563-a8f7-b32bba27d707',
    new Date('2021-12-12T17:34:00-05:00'),
    0,
    '["a","b","c"]',
  );

  expect(response.getKeyComponents()).toEqual([response.requestId]);
});

test('response.toObject()', () => {
  const response = new ServiceResponse(
    'ffbc9005-c62a-4563-a8f7-b32bba27d707',
    new Date('2021-12-12T17:34:00-05:00'),
    0,
    '["a","b","c"]',
  );
  const obj = {
    requestId: 'ffbc9005-c62a-4563-a8f7-b32bba27d707',
    time: new Date('2021-12-12T17:34:00-05:00'),
    statusCode: 0,
    returnValue: '["a","b","c"]',
  };

  expect(response.toObject()).toEqual(obj);
});

test('response.serialize()', () => {
  const response = new ServiceResponse(
    'ffbc9005-c62a-4563-a8f7-b32bba27d707',
    new Date('2021-12-12T17:34:00-05:00'),
    0,
    '["a","b","c"]',
  );
  const serialized =
    '{"requestId":"ffbc9005-c62a-4563-a8f7-b32bba27d707","time":"2021-12-12T22:34:00.000Z",' +
    '"statusCode":0,"returnValue":"[\\"a\\",\\"b\\",\\"c\\"]"}';

  expect(response.serialize()).toEqual(serialized);
});

test('response.validate()', () => {
  const response = new ServiceResponse('123456', new Date(0));

  expect(() => response.validate()).toThrow(/request ID/);
  response.requestId = 'ffbc9005-c62a-4563-a8f7-b32bba27d707';

  expect(() => response.validate()).toThrow(/response time/);
  response.time = new Date('2021-12-12T17:34:00-05:00');

  expect(() => response.validate()).not.toThrow();
});

test('ServiceResponse.fromObject()', () => {
  const obj = {
    requestId: 'ffbc9005-c62a-4563-a8f7-b32bba27d707',
    time: new Date('2021-12-12T17:34:00-05:00'),
    statusCode: 0,
    returnValue: '["a","b","c"]',
  };

  const response = ServiceResponse.fromObject(obj);
  expect(response.requestId).toEqual(obj.requestId);
  expect(response.time).toEqual(obj.time);
  expect(response.returnValue).toEqual(obj.returnValue);
  expect(response.statusCode).toEqual(obj.statusCode);
});

test('ServiceResponse.deserialize()', () => {
  const expected = new ServiceResponse(
    'ffbc9005-c62a-4563-a8f7-b32bba27d707',
    new Date('2021-12-12T17:34:00-05:00'),
    0,
    '["a","b","c"]',
  );
  const serialized =
    '{"requestId":"ffbc9005-c62a-4563-a8f7-b32bba27d707","time":"2021-12-12T17:34:00-05:00",' +
    '"statusCode":0,"returnValue":"[\\"a\\",\\"b\\",\\"c\\"]"}';

  const actual = ServiceResponse.deserialize(serialized);
  expect(actual.requestId).toEqual(expected.requestId);
  expect(actual.time).toEqual(expected.time);
  expect(actual.returnValue).toEqual(expected.returnValue);
  expect(actual.statusCode).toEqual(expected.statusCode);

  expect(() => ServiceResponse.deserialize('\x00')).toThrow();
});
