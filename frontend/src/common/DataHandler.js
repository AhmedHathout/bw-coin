/* ************************************************************************** */
/*                                                                            */
/*  DataHandler.js                                                            */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created:                                                 by elhmn        */
/*   Updated: Fri Feb 22 06:19:19 2019                        by elhmn        */
/*                                                                            */
/* ************************************************************************** */

export const JsonRequest = function
({
    url,
    method = 'GET',
    onsuccess,
    onerror,
    onstatus,
    data = null,
	isFormData = false,
})
{
	let obj = new XMLHttpRequest();
	data = isFormData ? data : data && JSON.stringify(data);

	obj.open(method, url, true);
	obj.onreadystatechange = function ()
	{
		if (obj.readyState === 4
            && (obj.status === 200
                || obj.status === 201))
		{
            onsuccess && onsuccess(obj.responseText);
		}
		else if (obj.readyState === 4
            && (obj.status === 400
                || obj.status === 403))
		{
			onerror && onerror(obj.responseText);
		}
        else if (obj.readyState === 4)
        {
            onstatus && onstatus(obj.responseText, obj.status);
        }
	};
	obj.send(data);
};

const DataHandler = {
    JsonRequest,
};

export default DataHandler;
