import request from './wrapper';

export class LocationService {
    getLocationName(id) {
        return request({
            method: 'get',
            url: `/locations/${id}`,
            data: {}
        });
    }
}
