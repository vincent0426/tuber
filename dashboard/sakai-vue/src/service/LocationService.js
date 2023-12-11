import request from './wrapper';

export class LocationService {
    getLocationName(id) {
        return request({
            method: 'get',
            url: `/locations/${id}`,
            data: {}
        });
    }
    driverSendLocation(id,lat,lng) {
        return request({
            method: 'get',
            url: `/ws/driver?trip_id=`+id,
            data: {"latitute": lat, "longitude":lng}

        });   
    }
    passengerGetLocation(id) {
        return request({
            method: 'get',
            url: `/ws/passenger?trip_id=`+id
        });   
    }
}
