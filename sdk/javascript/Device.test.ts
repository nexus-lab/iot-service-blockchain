import Device from './Device';

test('device.getKeyComponents()', () => {
  const device = new Device(
    'device1',
    'org1',
    'device1',
    'Device of Org1 User1',
    new Date('2021-12-12T17:34:00-05:00'),
  );

  expect(device.getKeyComponents()).toEqual([device.organizationId, device.id]);
});

test('device.toObject()', () => {
  const device = new Device(
    'device1',
    'org1',
    'device1',
    'Device of Org1 User1',
    new Date('2021-12-12T17:34:00-05:00'),
  );
  const obj = {
    id: 'device1',
    organizationId: 'org1',
    name: 'device1',
    description: 'Device of Org1 User1',
    lastUpdateTime: new Date('2021-12-12T17:34:00-05:00'),
  };

  expect(device.toObject()).toEqual(obj);
});

test('device.serialize()', () => {
  const device = new Device(
    'device1',
    'org1',
    'device1',
    'Device of Org1 User1',
    new Date('2021-12-12T17:34:00-05:00'),
  );
  const serialized =
    '{"id":"device1","organizationId":"org1","name":"device1",' +
    '"description":"Device of Org1 User1",' +
    '"lastUpdateTime":"2021-12-12T22:34:00.000Z"}';

  expect(device.serialize()).toEqual(serialized);
});

test('device.validate()', () => {
  const device = new Device('', '', '');

  expect(() => device.validate()).toThrow(/device ID/);
  device.id = 'device1';

  expect(() => device.validate()).toThrow(/organization ID/);
  device.organizationId = 'org1';

  expect(() => device.validate()).toThrow(/device name/);
  device.name = 'device1';

  expect(() => device.validate()).toThrow(/last update time/);
  device.lastUpdateTime = new Date('2021-12-12T17:34:00-05:00');

  expect(() => device.validate()).not.toThrow();
});

test('Device.fromObject()', () => {
  const obj = {
    id: 'device1',
    organizationId: 'org1',
    name: 'device1',
    description: 'Device of Org1 User1',
    lastUpdateTime: new Date('2021-12-12T17:34:00-05:00'),
  };

  const device = Device.fromObject(obj);
  expect(device.id).toEqual(obj.id);
  expect(device.organizationId).toEqual(obj.organizationId);
  expect(device.name).toEqual(obj.name);
  expect(device.description).toEqual(obj.description);
  expect(device.lastUpdateTime).toEqual(obj.lastUpdateTime);
});

test('Device.deserialize()', () => {
  const expected = new Device(
    'device1',
    'org1',
    'device1',
    'Device of Org1 User1',
    new Date('2021-12-12T17:34:00-05:00'),
  );
  const serialized =
    '{"id":"device1","organizationId":"org1","name":"device1",' +
    '"description":"Device of Org1 User1",' +
    '"lastUpdateTime":"2021-12-12T17:34:00-05:00"}';

  const actual = Device.deserialize(serialized);
  expect(actual.id).toEqual(expected.id);
  expect(actual.organizationId).toEqual(expected.organizationId);
  expect(actual.name).toEqual(expected.name);
  expect(actual.description).toEqual(expected.description);
  expect(actual.lastUpdateTime).toEqual(expected.lastUpdateTime);

  expect(() => Device.deserialize('\x00')).toThrow();
});
